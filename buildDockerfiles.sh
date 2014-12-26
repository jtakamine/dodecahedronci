#!/bin/bash

function process {
	cp $1 Dockerfile
	params=$@
	command=( "docker" "build" "${params[@]/$1}" "." )
	${command[@]}
}

if [ -z $1 ]
then
	#for dFile in *.Dockerfile
	#do
	#	process $dFile
	#done

	error "You must supply a path to the .Dockerfile to build"
else
	process $@ 
fi
