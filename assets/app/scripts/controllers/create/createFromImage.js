"use strict";

angular.module("openshiftConsole")
  .controller("CreateFromImageController", function ($scope,
      Logger,
      $q,
      $routeParams,
      DataService,
      ProjectsService,
      Navigate,
      ApplicationGenerator,
      TaskList,
      failureObjectNameFilter,
      $filter,
      $parse,
      SOURCE_URL_PATTERN
    ){
    var displayNameFilter = $filter('displayName');
    var humanize = $filter('humanize');

    $scope.projectName = $routeParams.project;
    $scope.sourceURLPattern = SOURCE_URL_PATTERN;

    ProjectsService
      .get($routeParams.project)
      .then(_.spread(function(project, context) {
        $scope.project = project;
        function initAndValidate(scope){

          if(!$routeParams.imageName){
            Navigate.toErrorPage("Cannot create from source: a base image was not specified");
          }
          if(!$routeParams.imageTag){
            Navigate.toErrorPage("Cannot create from source: a base image tag was not specified");
          }

          scope.emptyMessage = "Loading...";
          scope.imageName = $routeParams.imageName;
          scope.imageTag = $routeParams.imageTag;
          scope.namespace = $routeParams.namespace;
          scope.buildConfig = {
            buildOnSourceChange: true,
            buildOnImageChange: true,
            buildOnConfigChange: true,
            envVars : {
            }
          };
          scope.deploymentConfig = {
            deployOnNewImage: true,
            deployOnConfigChange: true,
            envVars : {
            }
          };
          scope.routing = {
            include: true,
            portOptions: []
          };
          scope.labels = {};
          scope.annotations = {};
          scope.scaling = {
            replicas: 1
          };

          scope.fillSampleRepo = function() {
            var annotations;
            if (!scope.image && !scope.image.metadata && !scope.image.metadata.annotations) {
              return;
            }

            annotations = scope.image.metadata.annotations;
            scope.buildConfig.sourceUrl = annotations.sampleRepo || "";
            scope.buildConfig.gitRef = annotations.sampleRef || "";
            scope.buildConfig.contextDir = annotations.sampleContextDir || "";
          };

          DataService.get("imagestreams", scope.imageName, {namespace: (scope.namespace || $routeParams.project)}).then(function(imageStream){
              scope.imageStream = imageStream;
              var imageName = scope.imageTag;
              DataService.get("imagestreamtags", imageStream.metadata.name + ":" + imageName, {namespace: scope.namespace}).then(function(imageStreamTag){
                  scope.image = imageStreamTag.image;
                  var env = $parse('dockerImageMetadata.ContainerConfig.Env')(imageStreamTag.image) || [];
                  angular.forEach(env, function(entry){
                    var pair = entry.split("=");
                    scope.deploymentConfig.envVars[pair[0]] = pair[1];
                  });

                  scope.routing.portOptions = ApplicationGenerator.parsePorts(imageStreamTag.image);
                  if (scope.routing.portOptions.length){
                    scope.routing.targetPort = scope.routing.portOptions[0];
                  } else {
                    scope.routing.include = false;
                  }
                }, function(){
                    Navigate.toErrorPage("Cannot create from source: the specified image could not be retrieved.");
                  }
                );
            },
            function(){
              Navigate.toErrorPage("Cannot create from source: the specified image could not be retrieved.");
            });
        }

        initAndValidate($scope);

        var ifResourcesDontExist = function(resources, namespace, scope){
          var result = $q.defer();
          var successResults = [];
          var failureResults = [];
          var remaining = resources.length;

          function _checkDone() {
            if (remaining === 0) {
              if(successResults.length > 0){
                //means some resources exist with the given nanme
                result.reject(successResults);
              }
              else
                //means no resources exist with the given nanme
                result.resolve(resources);
            }
          }

          resources.forEach(function(resource) {
            var resourceName = DataService.kindToResource(resource.kind);
            if (!resourceName) {
              failureResults.push({data: {message: "Unrecognized kind: " + resource.kind + "."}});
              remaining--;
              _checkDone();
              return;
            }
            DataService.get(resourceName, resource.metadata.name, {namespace: (namespace || $routeParams.project)}, {errorNotification: false}).then(
              function (data) {
                successResults.push(data);
                remaining--;
                _checkDone();
              },
              function (data) {
                failureResults.push(data);
                remaining--;
                _checkDone();
              }
            );
          });
          return result.promise;
        };

        var createResources = function(resources){
          var titles = {
            started: "Creating application " + $scope.name + " in project " + $scope.projectDisplayName(),
            success: "Created application " + $scope.name + " in project " + $scope.projectDisplayName(),
            failure: "Failed to create " + $scope.name + " in project " + $scope.projectDisplayName()
          };
          var helpLinks = {};

          TaskList.clear();
          TaskList.add(titles, helpLinks, function(){
            var d = $q.defer();
            DataService.createList(resources, context)
              //refactor these helpers to be common for 'newfromtemplate'
              .then(function(result) {
                    var alerts = [];
                    var hasErrors = false;
                    if (result.failure.length > 0) {
                      hasErrors = true;
                      result.failure.forEach(
                        function(failure) {
                          var objectName = failureObjectNameFilter(failure) || "object";
                          alerts.push({
                            type: "error",
                            message: "Cannot create " + humanize(objectName).toLowerCase() + ". ",
                            details: failure.data.message
                          });
                        }
                      );
                      result.success.forEach(
                        function(success) {
                          alerts.push({
                            type: "success",
                            message: "Created " + humanize(success.kind).toLowerCase() + " \"" + success.metadata.name + "\" successfully. "
                          });
                        }
                      );
                    } else {
                      alerts.push({ type: "success", message: "All resources for application " + $scope.name +
                        " were created successfully."});
                    }
                    d.resolve({alerts: alerts, hasErrors: hasErrors});
                  }
                );
                return d.promise;
              },
              function(result) { // failure
                $scope.alerts["create"] =
                  {
                    type: "error",
                    message: "An error occurred creating the application.",
                    details: "Status: " + result.status + ". " + result.data
                  };
              }
            );
          Navigate.toNextSteps($scope.name, $scope.projectName);
        };

        var elseShowWarning = function(){
          $scope.nameTaken = true;
          $scope.disableInputs = false;
        };

        $scope.projectDisplayName = function() {
          return displayNameFilter(this.project) || this.projectName;
        };

        $scope.createApp = function(){
          $scope.disableInputs = true;
          var resourceMap = ApplicationGenerator.generate($scope);
          //init tasks
          var resources = [];
          angular.forEach(resourceMap, function(value, key){
            if(value !== null){
              Logger.debug("Generated resource definition:", value);
              resources.push(value);
            }
          });

          ifResourcesDontExist(resources, $scope.projectName, $scope)
            .then(createResources, elseShowWarning);
        };
      }));
  });

