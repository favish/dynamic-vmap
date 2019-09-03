# Generates Dynamic VMap files based on input parameters
This is a very simple proof of concept but should scale without too much trouble with the addition of a k8s service, pod, and ingress that specifies an nginx and SSL (lets encrypt) annotation. 

You may want to use the released build artifact (binary) of this project in the creation of a new project that specifies these things as it is poutside the scope of this generally useful too that we can make publicly available.   

## TODO
* Use environment vars to specify the VAST tags instead of hard code to allow us to use in mulitple places
* Break monofunction into something more comprehensible
* Add optional "duration" query parameter for the video length to automatically space out ads.
* Eventually create a map of minimum video lengths before new ad formats take effect (e.g. all video longer than 0 seconds get only preroll, all videos longer than 60 second get preroll and post roll, all longer than 180 get...)
