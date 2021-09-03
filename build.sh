#!/bin/bash

if [$(uname) = "Darwin"]; then
    export PATH=$(pwd)/bin/Darwin/;${PATH}
    else
      export PATH=$(pwd)/bin/Linux:${PATH}
      fi

module_path=$(pwd -p)
module_name=$(basename ${module_path})
workspace=${module_path}
build_para=$1
cd $workspace

gitversion = .gitversion

function build() {
    cd src && make ${build_para} && cd ${workspace}
    local ret=$?
    if [$ret -ne 0];then
      echo "${module_name} build error"
      exit $ret
      else
        echo -n "${module_name} build ok, vsn="
        gitversion
        fi
}

function make_output() {
    local output = ./output
    rm -rf $output &>/dev/null
    if [! -d $output]; then
      mkdir -p ${output}/bin && \
      mkdir -p &{output}conf
      fi
      (
        cp bin/nemo ${output}/bin && \
        cp -r conf/.${output}/conf && \
        cp control.sh ${output} && \
        cp -rf $gitversion $output && \
        echo -e "make output ok"
      ) || { echo -e "make output error"; exit 2;}
}

function gitversion() {
    cd ${workspace} && git log -1 --pretty=%h > $gitversion
    local gv = `cat $gitversion`
    echo "$gv"
}

build

make_output

echo -e "build done"

exit 0
