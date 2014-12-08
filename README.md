DodecahedronCI
==============

Required Environment Variables:
* DODEC\_HOME
* DODEC\_GITHUB\_USER
* DODEC\_GITHUB\_PASSWORD
* DODEC\_DOCKER\_USER
* DODEC\_DOCKER\_PASSWORD
* DODEC\_DOCKER\_EMAIL

TODO:
* Make things work well when containerized
* Make sure permissions are correctly set out of the box (file permissions/execution permissions/etc.)
* Use environment variables to set configuration:
  * Github credentials
  * Docker credentials
  * Top-level directory where dodecci can clone repositories
* Avoid having to log in when pushing to Docker repo
* Map Dockerfiles to Docker image names + Docker credentials (to log into repo)
* Filter by "ends with Dockerfile" to search for Dockerfiles to build (instead of "equals Dockerfile")
