#!/bin/bash
ibmcloud login --apikey $APIKEY -a https://cloud.ibm.com -r us-south

ibmcloud target -o advowork@us.ibm.com -s dev

ibmcloud account list