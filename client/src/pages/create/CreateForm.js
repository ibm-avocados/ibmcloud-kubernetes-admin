import React from "react";

import NotificationEmail from "./NotificationEmail";

import {
  Button,
  Column,
  DatePicker,
  DatePickerInput,
  Dropdown,
  Form,
  FormLabel,
  Grid,
  InlineLoading,
  ModalWrapper,
  MultiSelect,
  RadioTile,
  Row,
  SelectItem,
  TextInput,
  TimePicker,
  TimePickerSelect,
  ToastNotification,
  Tooltip,
} from "carbon-components-react";

import geos from "../../common/data/geo";

import styles from "./CreateForm.module.css";

import "./CreateForm.css";
import WorkshopAccount from "./WorkshopAccount";

const Spacer = ({ height }) => <div style={{ marginTop: height }} />;

const Divider = ({ width }) => <div style={{ marginRight: width }} />;

const grab = async (url, options, retryCount = 0) => {
  const response = await fetch(url, options);
  if (response.status !== 200) {
    if (retryCount > 0) {
      console.log("failure in request. retrying again");
      return await grab(url, options, retryCount - 1);
    }
    throw Error(data);
  }
  const data = await response.json();
  return data;
};

const CreateForm = ({ accountID }) => {
  // radio tile
  const [kubernetesSelected, setKubernetesSelected] = React.useState(true);
  const [openshiftSelected, setOpenshiftSelected] = React.useState(false);

  // values
  const [kuberntesVersions, setKubernetesVersions] = React.useState([]);
  const [openshiftVersions, setOpenshiftVersions] = React.useState([]);
  const [workerZones, setWorkerZones] = React.useState([]);
  // const [privateVlans, setPrivateVlans] = React.useState([]);
  // const [publicVlans, setPublicVlans] = React.useState([]);
  const [clusterNamePrefix, setClusterNamePrefix] = React.useState("");
  const [clusterCount, setClusterCount] = React.useState("1");
  const [workerCount, setWorkerCount] = React.useState("1");
  const [tags, setTags] = React.useState("");
  const [flavors, setFlavors] = React.useState([]);
  const [resourceGroups, setResourceGroups] = React.useState([]);
  // selected values
  const [selectedKubernetes, setSelectedKuberetes] = React.useState(null);
  const [selectedOpenshift, setSelectedOpenshift] = React.useState(null);
  const [selectedRegion, setSelectedRegion] = React.useState(null);
  const [selectedWorkerZones, setSelectedWorkerZones] = React.useState([]);
  // const [selectedWorkerZone, setSelectedWorkerZone] = React.useState(null);
  // const [selectedPrivateVlan, setSelecetedPrivateVlan] = React.useState(null);
  // const [selectedPublicVlan, setSelecetedPublicVlan] = React.useState(null);
  const [zoneClusterCount, setZoneClusterCount] = React.useState(null);
  const [selectedFlavor, setSelectedFlavor] = React.useState(null);
  const [selectedGroup, setSelectedGroup] = React.useState(null);
  //scheduling helpers
  const [startTimeAMPM, setStartTimeAMPM] = React.useState("AM");
  const [endTimeAMPM, setEndTimeAMPM] = React.useState("AM");
  const [apiKey, setApiKey] = React.useState("");
  const [apiKeyValid, setApiKeyValid] = React.useState(false);
  const [startTime, setStartTime] = React.useState("");
  const [endTime, setEndTime] = React.useState("");
  const [dateRange, setDateRange] = React.useState([]);

  const [isWorkshop, setIsWorkshop] = React.useState(false);
  const [githubIssue, setGithubIssue] = React.useState("");
  const [userPerCluster, setUserPerCluster] = React.useState("1");

  // ui indicators
  const [creating, setCreating] = React.useState(false);
  const [loaderDescription, setLoaderDescription] = React.useState("");
  const [createSuccess, setCreateSuccess] = React.useState(false);
  const [scheduleSuccess, setScheduleSuccess] = React.useState(false);
  const [toast, setToast] = React.useState(null);

  // notification specific states

  const [selectedEmails, setSelectedEmails] = React.useState([]);
  const [awxWorkflowJobTemplates, setAWXWorkflowJobTemplates] = React.useState([]);

  React.useEffect(() => {
    const loadWorkflowJobTemplate = async() => {
      try {
        const templates = await grab("/api/v1/awx/workflowjobtemplate?labels=kubernetes")
        console.log(templates);
        setAWXWorkflowJobTemplates(templates)
      }
      catch (e) {
        console.log(e);
      }
    }
    loadWorkflowJobTemplate()

    const loadVersions = async () => {
      try {
        const versions = await grab("/api/v1/clusters/versions");
        if (versions) {
          setKubernetesVersions(versions.kubernetes);
          setOpenshiftVersions(versions.openshift);
        }
      } catch (e) {
        console.log(e);
      }
    };
    loadVersions();

    const loadResourceGroups = async () => {
      try {
        const resourceGroups = await grab(
          `/api/v1/resourcegroups/${accountID}`
        );
        if (resourceGroups) {
          setResourceGroups(resourceGroups.resources);
          console.log(resourceGroups.resources);
        }
      } catch (e) {
        console.log(e);
      }
    };

    loadResourceGroups();

    const checkAPIKey = async () => {
      try {
        const apiKey = await fetch("/api/v1/schedule/api", {
          method: "post",
          body: JSON.stringify({
            accountID: accountID,
          }),
        });
        if (apiKey.status === 200) {
          setApiKeyValid(true);
          setApiKey("your-api-key-will-be-pulled-from-db");
        }
      } catch (e) {
        console.log(e);
      }
    };
    checkAPIKey();
    getWorkerZoneClusterCount();
  }, [accountID]);

  const resetState = () => {
    setTags("");
    setWorkerCount(1);
    setClusterCount(1);
    setClusterNamePrefix("");
    setStartTime("");
    setEndTime("");
    setDateRange([]);
    setSelectedKuberetes(null);
    setSelectedOpenshift(null);
    setSelectedRegion(null);
    // setSelectedWorkerZone(null);
    setSelectedWorkerZones([]);
    // setSelecetedPrivateVlan(null);
    // setSelecetedPublicVlan(null);
    setSelectedFlavor(null);
    setSelectedGroup(null);
  };

  const toggleRadio = () => {
    setKubernetesSelected(!kubernetesSelected);
    setOpenshiftSelected(!openshiftSelected);
    setFlavorOnClusterType(flavors, kubernetesSelected);
  };

  const getWorkerZones = async (geo) => {
    try {
      const locations = await grab(`/api/v1/clusters/${geo}/locations`, {
        Method: "GET",
      });
      if (locations && locations.length > 0) {
        setWorkerZones(locations);
      }
    } catch (e) {
      console.log(e);
    }
  };

  const getWorkerZoneClusterCount = async () => {
    try {
      const zoneCount = await grab(`/api/v1/clusters/locations/info`, {
        Method: "GET",
      });
      if(zoneCount) {
        console.log(zoneCount);
        setZoneClusterCount(zoneCount);
      }
    } catch (e) {
      console.log(e);
    }
  }

  const getVlan = async (datacenter) => {
    try {
      const vlans = await grab(`/api/v1/clusters/${datacenter}/vlans`);
      const privateVlans = vlans.filter((vlan) => vlan.type === "private");
      const publicVlans = vlans.filter((vlan) => vlan.type === "public");
      if (privateVlans === null && publicVlans === null) {
        return [];
      } else if (privateVlans === null || publicVlans === null) {
        throw new Error("need both pair");
      }

      if (privateVlans.length > 0 && publicVlans.length > 0) {
        return findMatchingVlans(privateVlans, publicVlans);
      }
      return [];
    } catch (e) {
      throw new Error(e);
    }
  };

  function getRandomInt(max) {
    return Math.floor(Math.random() * Math.floor(max));
  }

  const findMatchingVlans = (privateVlans, publicVlans) => {
    let pairs = [];

    for (let privateVlan of privateVlans) {
      console.log(privateVlan);
      const privateMatch = privateVlan.properties.primary_router.substring(1);
      for (let publicVlan of publicVlans) {
        const publicMatch = publicVlan.properties.primary_router.substring(1);
        if (privateMatch === publicMatch) {
          pairs.push([privateVlan.id, publicVlan.id]);
        }
      }
    }

    if (pairs.length === 0) {
      throw new Error("no matching vlan found");
    }

    return pairs[getRandomInt(pairs.length)];
  };

  // const getVlans = async (datacenter) => {
  //   try {
  //     const vlans = await grab(`/api/v1/clusters/${datacenter}/vlans`);
  //     const privateVlans = vlans.filter((vlan) => vlan.type === "private");
  //     if (privateVlans && privateVlans.length > 0) {
  //       setPrivateVlans(privateVlans);
  //       setSelecetedPrivateVlan(privateVlans[0]);
  //     }
  //     const publicVlans = vlans.filter((vlan) => vlan.type === "public");
  //     if (publicVlans && publicVlans.length > 0) {
  //       setPublicVlans(publicVlans);
  //       setSelecetedPublicVlan(publicVlans[0]);
  //     }
  //   } catch (e) {
  //     console.log(e);
  //   }
  // };

  const getFlavors = async (datacenter) => {
    try {
      const flav = await grab(
        `/api/v1/clusters/${datacenter}/machine-types?type=virtual&os=UBUNTU_18_64&cpuLimit=8&memoryLimit=32`
      );
      if (flav) {
        setFlavorOnClusterType(flav, openshiftSelected);
      }
    } catch (e) {
      console.log(e);
    }
  };

  const setFlavorOnClusterType = (flav, isOpenshift) => {
    if (!isOpenshift) {
      setFlavors(flav);
    } else {
      setFlavors(flav.filter((flavor) => Number(flavor.cores) > 2));
    }
  };

  const getVersionString = (versions, version) => {
    const index = versions.indexOf(version);
    let substring = "stable";
    if (index === versions.length - 1) {
      substring = "latest";
    } else if (version.default) {
      substring = "stable, default";
    }

    return `${version.major}.${version.minor}.${version.patch} (${substring})`;
  };

  const getKuberntesVersionString = (version) =>
    getVersionString(kuberntesVersions, version);

  const getOpenshiftVersionString = (version) =>
    getVersionString(openshiftVersions, version);

  const getVlanString = (vlan) =>
    `${vlan.id}-${vlan.properties.vlan_number}-${vlan.properties.primary_router}`;

  const onGeoSelected = (geo) => {
    getWorkerZones(geo.id);
    setSelectedRegion(geo);
    setSelectedWorkerZones([]);
    // setSelectedWorkerZone(null);
    // setPrivateVlans([]);
    // setSelecetedPrivateVlan(null);
    // setPublicVlans([]);
    // setSelecetedPublicVlan(null);
  };

  const onWorkerZonesSelected = ({ selectedItems }) => {
    setSelectedWorkerZones(selectedItems);
    console.log(selectedItems);
    if (selectedItems.length > 0) {
      getFlavors(selectedItems[0].id);
    }
  };

  // const onWorkerZoneSelected = (zone) => {
  //   // setSelectedWorkerZone(zone);
  //   // setPrivateVlans([]);
  //   // setSelecetedPrivateVlan(null);
  //   // setPublicVlans([]);
  //   // setSelecetedPublicVlan(null);
  //   // getVlans(zone.id);
  //   getFlavors(zone.id);
  // };

  const validTag = (tags) => {
    const re = /^[A-Za-z,0-9:_ .-]+$/;
    const valid = re.test(tags);
    return !valid;
  };
  const numToStr = (num) => {
    const numstr = num.toString();
    const pad = "000";
    return pad.substring(0, pad.length - numstr.length) + numstr;
  };

  const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

  const getCreateRequest = () => {
    let version = "";
    if (kubernetesSelected) {
      const { major, minor, patch } = selectedKubernetes;
      version = `${major}.${minor}.${patch}`;
    } else {
      const { major, minor } = selectedOpenshift;
      version = `${major}.${minor}_openshift`;
    }

    let defaultWorkerPoolEntitlement = "";
    if (openshiftSelected) {
      defaultWorkerPoolEntitlement = "cloud_pak";
    }

    const range = Number(clusterCount);
    let request = [];
    for (let i = 1; i <= range; i++) {
      const suffix = numToStr(i);
      const name = `${clusterNamePrefix}-${suffix}`;

      const ClusterRequest = {
        name,
        prefix: "",
        skipPermPrecheck: false,
        // dataCenter: selectedWorkerZone.id,
        defaultWorkerPoolName: "",
        defaultWorkerPoolEntitlement,
        disableAutoUpdate: true,
        noSubnnet: false,
        podSubnet: "",
        serviceSubnet: "",
        machineType: selectedFlavor.name,
        // privateVlan: selectedPrivateVlan.id,
        // publicVlan: selectedPublicVlan.id,
        masterVersion: version,
        workerNum: Number(workerCount),
        diskEncryption: true,
        isolation: "public",
        gatewayEnabled: false,
        privateSeviceEndpoint: false,
        publicServiceEndpoint: false,
      };

      const CreateClusterRequest = {
        clusterRequest: ClusterRequest,
        resourceGroup: selectedGroup.id,
      };
      request.push(CreateClusterRequest);
    }
    return request;
  };

  const onCreateClicked = async () => {
    console.log("creating clusters");
    setCreating(true);
    setCreateSuccess(false);

    const request = getCreateRequest();
    const range = request.length;
    let errors = [];
    for (let i = 0; i < range; i++) {
      setLoaderDescription(`Creating Cluster ${i + 1} of ${range}`);
      console.log("creating cluster ", i);

      const CreateClusterRequest = request[i];

      try {
        const workerZone =
          selectedWorkerZones[i % selectedWorkerZones.length].id;
        setLoaderDescription(`Getting Vlan for Cluster ${i + 1} of ${range}`);
        const vlan = await getVlan(workerZone);

        // if there is a vlan set it
        // else set empty
        // on ibmcloud if a datacenter does not have a vlan
        // it will be created
        if (vlan.length !== 0) {
          CreateClusterRequest.clusterRequest.privateVlan = vlan[0];
          CreateClusterRequest.clusterRequest.publicVlan = vlan[1];
        }

        CreateClusterRequest.clusterRequest.dataCenter = workerZone;

        console.log(
          `Worker Zone: ${CreateClusterRequest.clusterRequest.workerZone},\tPrivate Vlan: ${CreateClusterRequest.clusterRequest.privateVlan},\tPublic Vlan: ${CreateClusterRequest.clusterRequest.publicVlan}`
        );

        console.log(CreateClusterRequest);

        const clusterResponse = await grab(
          "/api/v1/clusters",
          {
            method: "post",
            body: JSON.stringify(CreateClusterRequest),
          },
          3
        );

        console.log("cluster created with id : ", clusterResponse.id);

        console.log("Sleeping 3s before trying to set tags");
        setLoaderDescription(`Preparing to Tag Cluster ${i + 1} of ${range}`);
        await sleep(15000);
        setLoaderDescription(`Tagging Cluster ${i + 1} of ${range}`);

        // comma separated tags.
        const tagPromises = tags.split(",").map(async (tag) => {
          try {
            await sleep(3000);
            const tagRequest = await grab(
              `/api/v1/clusters/${clusterResponse.id}/settag`,
              {
                method: "post",
                body: JSON.stringify({
                  tag,
                  resourceGroup: selectedGroup.id,
                }),
              },
              3
            );
            return tagRequest;
          } catch (e) {
            return undefined;
          }
        });
        const result = await Promise.all(tagPromises);
        console.log(result);
      } catch (e) {
        errors.push(e);
        console.log("Error creating cluster", e);
      }
    }

    const datacenters = selectedWorkerZones.map(v => v.id).join(", ");

    setToast({
      title: "Cluster Created",
      subtitle: `${clusterCount} ${
        kubernetesSelected ? "Kubernetes" : "Openshift"
      } Cluster Creation Attempted. ${errors.length} Error`,
      kind: errors.length === 0 ? "success" : "error",
      caption: `Datacenter(s): ${datacenters}`,
    });
    setCreateSuccess(true);
    setCreating(false);
    resetState();
  };

  const shouldCreateBeDisabled = () => {
    let versionSelected = false;
    if (kubernetesSelected) {
      versionSelected = !!selectedKubernetes;
    } else {
      versionSelected = !!selectedOpenshift;
    }

    const groupSelected = !!selectedGroup;
    const geoSelected = !!selectedRegion;
    const zoneSelected = selectedWorkerZones.length > 0;
    const flavorSelected = !!selectedFlavor;

    const hasClusterCount = clusterCount && clusterCount !== "";
    const hasWorkerCount = workerCount && workerCount !== "";
    const hasNamePrefix = clusterNamePrefix && clusterNamePrefix !== "";
    const hasTags = tags && tags !== "";

    return !(
      versionSelected &&
      groupSelected &&
      geoSelected &&
      zoneSelected &&
      flavorSelected &&
      hasClusterCount &&
      hasWorkerCount &&
      hasNamePrefix &&
      hasTags
    );
  };

  const shouldSchedulingBeDisabled = () => {
    return true;
    //TODO: scheduler is broken right now. enable once fixed/
    //return shouldCreateBeDisabled();
  };

  const shouldScheduleSubmitBeDisabled = () => {
    const dateSet = dateRange.length >= 1;
    const createTimetSet = !timeInvalid(startTime);
    const endTimeSet = !timeInvalid(endTime);

    return !(dateSet && createTimetSet && endTimeSet && apiKeyValid);
  };

  const onScheduleSubmit = async () => {
    if (shouldScheduleSubmitBeDisabled()) {
      return false;
    }

    setScheduleSuccess(false);
    let version = "";
    if (kubernetesSelected) {
      const { major, minor, patch } = selectedKubernetes;
      version = `${major}.${minor}.${patch}`;
    } else {
      const { major, minor } = selectedOpenshift;
      version = `${major}.${minor}_openshift`;
    }

    let defaultWorkerPoolEntitlement = "";
    if (openshiftSelected) {
      defaultWorkerPoolEntitlement = "cloud_pak";
    }
    const name = `${clusterNamePrefix}`;

    const workerZones = selectedWorkerZones.map(v => v.id);

    const ScheduleRequest = {
      name: name,
      prefix: "",
      skipPermPrecheck: false,
      dataCenters: workerZones,
      defaultWorkerPoolName: "",
      defaultWorkerPoolEntitlement,
      disableAutoUpdate: true,
      noSubnnet: false,
      podSubnet: "",
      serviceSubnet: "",
      machineType: selectedFlavor.name,
      masterVersion: version,
      workerNum: Number(workerCount),
      diskEncryption: true,
      isolation: "public",
      GatewayEnabled: false,
      privateSeviceEndpoint: false,
      publicServiceEndpoint: false,
    };

    const ScheduleClusterRequest = {
      scheduleRequest: ScheduleRequest,
      resourceGroup: selectedGroup.id,
    };

    let startDate = dateRange[0];
    let endDate = dateRange[1];
    let startHour = Number(startTime.split(":")[0]);
    startHour += startTimeAMPM === "PM" ? 12 : 0;
    const startMinute = Number(startTime.split(":")[1]);
    let endHour = Number(endTime.split(":")[0]);
    endHour += endTimeAMPM === "PM" ? 12 : 0;
    const endMinute = Number(endTime.split(":")[1]);

    startDate.setTime(
      startDate.getTime() + startHour * 60 * 60 * 1000 + startMinute * 60 * 1000
    );
    const createAt = startDate.getTime() / 1000;

    endDate.setTime(
      endDate.getTime() + endHour * 60 * 60 * 1000 + endMinute * 60 * 1000
    );
    const destroyAt = endDate.getTime() / 1000;

    const password = kubernetesSelected ? "ikslab" : "oslab";

    const schedule = {
      createAt: createAt,
      destroyAt: destroyAt,
      status: "scheduled",
      tags: tags,
      count: clusterCount,
      userCount: userPerCluster,
      scheduleRequest: ScheduleClusterRequest,
      clusters: [],
      notifyEmails: selectedEmails,
      eventName: clusterNamePrefix,
      password: password,
      resourceGroupName: selectedGroup.name,
      githubIssueNumber: githubIssue,
      isWorkshop: isWorkshop,
    };

    /*
    	EventName         string               `json:"eventName"`
	Password          string               `json:"password"`
	ResourceGroupName string               `json:"resourceGroupName"`
    */

    try {
      const response = await grab(`/api/v1/schedule/${accountID}/create`, {
        method: "post",
        body: JSON.stringify(schedule),
      });
      console.log("schedule set");
    } catch (e) {
      console.log(e);
    }
    setToast({
      title: "Cluster Scheduled",
      subtitle: `${clusterCount} ${
        kubernetesSelected ? "Kubernetes" : "Openshift"
      } Clusters Scheduled`,
      kind: "success",
      caption: `Create At : ${startDate.toLocaleString()} Delete At : ${endDate.toLocaleString()}`,
    });
    setScheduleSuccess(true);
    resetState();
    return true;
  };

  const onSubmitAPIKeyClicked = async () => {
    try {
      const response = await grab("/api/v1/schedule/api/create", {
        method: "post",
        body: JSON.stringify({
          accountID: accountID,
          apiKey: apiKey,
        }),
      });
      console.log(response);
      setApiKeyValid(true);
      setApiKey("your-api-key-will-be-pulled-from-db");
    } catch (e) {
      console.log(e);
    }
  };

  const onDeleteAPIKeyClicked = async () => {
    try {
      const response = await grab("/api/v1/schedule/api", {
        method: "delete",
        body: JSON.stringify({
          accountID: accountID,
        }),
      });
      console.log("api key removed");
      setApiKey("");
      setApiKeyValid(false);
    } catch (e) {
      console.log(e);
    }
  };

  const timeInvalid = (time) => {
    const re = /^(0[0-9]|1[0-2]):[0-5][0-9]$/;
    return !re.test(time);
  };

  const renderFlavors = (item) => {
    if (item) {
      return (
        <div style={{ position: "absolute" }}>
          <p>
            {item.cores} vCPUs
            {item.memory} RAM
          </p>
          <p>{item.name}</p>
        </div>
      );
    }
    return null;
  };

  const setMinDate = () => {
    let d = new Date();
    return `${d.getMonth() + 1}/${d.getDate()}/${d.getFullYear}`;
  };

  const setMaxDate = () => {
    let d = new Date();
    d.setDate(d.getDate() + 60);
    return `${d.getMonth() + 1}/${d.getDate()}/${d.getFullYear}`;
  };

  const getZoneText = (zone) => {
    // (zone) => (zone ? zone.id + ` (${zoneClusterCount[zone.id]?zoneClusterCount[zone.id]:0})` : "")
    return zone && zoneClusterCount ? `${zone.id} (${zoneClusterCount[zone.id]?zoneClusterCount[zone.id]:0})`:"";
  }

  return (
    <>
      <Form>
        <Grid>
          <FormLabel>
            <Tooltip triggerText="Cluster type and version">
              The container platform type and version for the cluster. Choose
              Kubernetes for a native Kubernetes experience on Ubuntu, or
              OpenShift to deliver Kubernetes apps on Red Hat Enterprise Linux.
            </Tooltip>
          </FormLabel>
          <Row>
            <Column md={4} lg={3}>
              <RadioTile
                disabled={creating}
                value="k8s"
                onClick={toggleRadio}
                checked={kubernetesSelected}
                className={styles.radio_tile}
              >
                <div className="radio-tile-content">
                  <img
                    alt="kubernetes logo"
                    className="radio-tile-image"
                    src="https://i.ibb.co/cDqxKBd/download.png"
                    height={100}
                    width={100}
                  />
                  <p>Kubernetes</p>
                  <Spacer height="16px" />
                </div>
                <Dropdown
                  id="kubernetes_version"
                  className="create-page-dropdown"
                  disabled={!kubernetesSelected || creating}
                  label="Select Version"
                  items={kuberntesVersions}
                  onChange={({ selectedItem }) =>
                    setSelectedKuberetes(selectedItem)
                  }
                  selectedItem={selectedKubernetes}
                  itemToString={(version) =>
                    version ? getKuberntesVersionString(version) : ""
                  }
                />
              </RadioTile>
            </Column>
            <Column md={4} lg={3}>
              <RadioTile
                disabled={creating}
                value="openshift"
                onClick={toggleRadio}
                checked={openshiftSelected}
                className={styles.radio_tile}
              >
                <div className="radio-tile-content">
                  <img
                    alt="openshift logo"
                    className="radio-tile-image"
                    src="https://i.ibb.co/0fFQCD2/openshift.png"
                    height={100}
                    width={100}
                  />
                  <p>OpenShift</p>
                  <Spacer height="16px" />
                </div>
                <Dropdown
                  id="openshift_version"
                  className="create-page-dropdown"
                  disabled={!openshiftSelected || creating}
                  label="Select Version"
                  items={openshiftVersions}
                  onChange={({ selectedItem }) =>
                    setSelectedOpenshift(selectedItem)
                  }
                  selectedItem={selectedOpenshift}
                  itemToString={(version) =>
                    version ? getOpenshiftVersionString(version) : ""
                  }
                />
              </RadioTile>
            </Column>
          </Row>
          <Spacer height="16px" />

          <Spacer height="16px" />
          <Row>
            <Column>
              <h2>Location</h2>
            </Column>
          </Row>
          <Spacer height="16px" />

          <FormLabel>Geography</FormLabel>

          <Row>
            <Column lg={6}>
              <Dropdown
                id="geo_select"
                className="create-page-dropdown"
                label="Select geo"
                items={geos}
                disabled={creating}
                selectedItem={selectedRegion}
                itemToString={(geo) => (geo ? geo.display_name : "")}
                onChange={({ selectedItem }) => onGeoSelected(selectedItem)}
              />
            </Column>
          </Row>
          <Spacer height="16px" />

          <Row>
            <Column lg={6}>
              <FormLabel>
                <Tooltip triggerText="Worker zone">
                  The data center where your worker pool will be located.
                </Tooltip>
              </FormLabel>
              <MultiSelect
                id="workerzones-select"
                itemToString={zone => getZoneText(zone)}
                items={workerZones}
                disabled={workerZones.length <= 0 || !zoneClusterCount || creating}
                className="create-page-multiselect"
                label="Select worker zone"
                onChange={(selected) => onWorkerZonesSelected(selected)}
              />
            </Column>
          </Row>

          {/* <Row>
            <Column lg={6}>
              <FormLabel>
                <Tooltip triggerText="Worker zone">
                  The data center where your worker pool will be located.
                </Tooltip>
              </FormLabel>
              <Dropdown
                id="worker_zone_select"
                className="create-page-dropdown"
                label="Select worker zone"
                itemToString={(zone) => (zone ? zone.id : "")}
                items={workerZones}
                selectedItem={selectedWorkerZone}
                disabled={workerZones.length <= 0 || creating}
                onChange={({ selectedItem }) =>
                  onWorkerZoneSelected(selectedItem)
                }
              />
            </Column>

            <Column md={4} lg={3}>
              <FormLabel>
                <Tooltip triggerText="Public vlan">
                  Allow your worker nodes to securely communicate to the
                  IBM-managed master through a virtual network. To expose your
                  apps to the public, configure external networking.
                </Tooltip>
              </FormLabel>
              <Dropdown
                id="public_vlan"
                className="create-page-dropdown"
                label="Select public vlan"
                disabled={publicVlans.length <= 0 || creating}
                items={publicVlans}
                itemToString={(vlan) => (vlan ? getVlanString(vlan) : "")}
                selectedItem={selectedPublicVlan}
                onChange={({ selectedItem }) =>
                  setSelecetedPublicVlan(selectedItem)
                }
              />
            </Column>

            <Column md={4} lg={3}>
              <FormLabel>
                <Tooltip triggerText="Private vlan">
                  Virtual network that allows private communication between
                  worker nodes in this cluster.
                </Tooltip>
              </FormLabel>
              <Dropdown
                id="private_vlan"
                className="create-page-dropdown"
                label="Select private vlan"
                disabled={privateVlans.length <= 0 || creating}
                items={privateVlans}
                itemToString={(vlan) => (vlan ? getVlanString(vlan) : "")}
                selectedItem={selectedPrivateVlan}
                onChange={({ selectedItem }) =>
                  setSelecetedPrivateVlan(selectedItem)
                }
              />
            </Column>
          </Row> */}
          <Spacer height="16px" />

          <Spacer height="16px" />
          <Row>
            <Column>
              <h2>Cluster Metadata</h2>
            </Column>
          </Row>
          <Spacer height="16px" />

          <Row>
            <Column md={4} lg={3}>
              <FormLabel>
                <Tooltip triggerText="Cluster name prefix">
                  Cluster name will be generated as prefix-001
                </Tooltip>
              </FormLabel>
              <TextInput
                value={clusterNamePrefix}
                onChange={(e) => setClusterNamePrefix(e.target.value.trim())}
                labelText=""
                disabled={creating}
                id="cluster_name"
                placeholder="mycluster"
              />
            </Column>

            <Column md={4} lg={3}>
              <TextInput
                value={clusterCount}
                onChange={(e) => setClusterCount(e.target.value.trim())}
                labelText="Cluster count"
                id="cluster_count"
                disabled={creating}
                placeholder="20"
                invalid={isNaN(clusterCount) || clusterCount === ""}
                invalidText="Should be a positive number"
              />
            </Column>

            <Column lg={3}>
              <FormLabel>Resouce group</FormLabel>
              <Dropdown
                id="resource_group"
                className="create-page-dropdown"
                label="Select resource group"
                items={resourceGroups}
                disabled={creating}
                itemToString={(item) => (item ? item.name : "")}
                selectedItem={selectedGroup}
                onChange={({ selectedItem }) => {
                  setSelectedGroup(selectedItem);
                }}
              />
            </Column>
          </Row>
          <Spacer height="16px" />

          <Row>
            <Column>
              <FormLabel>
                <Tooltip triggerText="Tags">
                  Use tags that would uniquely identify this set of clusters.
                </Tooltip>
              </FormLabel>
              <TextInput
                labelText=""
                id="tag_text"
                className="tag-text-input"
                placeholder="tag1, tag2, tag3"
                value={tags}
                disabled={creating}
                onChange={(e) => setTags(e.target.value)}
                invalid={tags !== "" && validTag(tags)}
                invalidText="valid tag is in the regex form ^[A-Za-z0-9:_ .-]+$"
              />
            </Column>
          </Row>
          <Spacer height="16px" />

          <Spacer height="16px" />
          <Row>
            <Column>
              <h2>Default worker pool</h2>
            </Column>
          </Row>
          <Spacer height="16px" />

          <Row>
            <Column md={4} lg={3}>
              <FormLabel>Worker nodes</FormLabel>
              <TextInput
                value={workerCount}
                onChange={(e) => setWorkerCount(e.target.value.trim())}
                labelText=""
                id="worker_nodes"
                disabled={creating}
                placeholder="1"
                invalid={isNaN(workerCount) || workerCount === ""}
                invalidText="Should be a positive number"
              />
            </Column>
          </Row>
          <Spacer height="16px" />

          <Row>
            <Column lg={6}>
              <FormLabel>
                <Tooltip triggerText="Flavor">
                  The amount of memory, cpu, and disk space allocated to each
                  worker node.
                </Tooltip>
              </FormLabel>
              <Dropdown
                id="machine_flavor"
                className="create-page-dropdown machine-flavor"
                label="Select Flavor"
                items={flavors}
                disabled={flavors.length <= 0 || creating}
                selectedItem={selectedFlavor}
                onChange={({ selectedItem }) => setSelectedFlavor(selectedItem)}
                itemToString={(item) => (item ? item.name : "")}
                itemToElement={(item) => renderFlavors(item)}
              />
            </Column>
          </Row>
          <Spacer height="16px" />

          {createSuccess || scheduleSuccess ? (
            <>
              <Spacer height="16px" />
              <ToastNotification
                title={toast.title}
                subtitle={toast.subtitle}
                kind={toast.kind}
                caption={toast.caption}
                timeout={0}
                style={{
                  minWidth: "50rem",
                  marginBottom: ".5rem",
                }}
              />
              <Spacer height="16px" />
            </>
          ) : (
            <></>
          )}

          <div>{awxWorkflowJobTemplates && awxWorkflowJobTemplates.map((v, i) => (<p key={i}>{v.name}</p>))}</div>

          <Spacer height="16px" />
          <Row>
            <Column>
              <div style={{ display: "flex", width: "500px" }}>
                {creating ? (
                  <InlineLoading
                    style={{ width: "250px" }}
                    description={loaderDescription}
                    status={createSuccess ? "finished" : "active"}
                  />
                ) : (
                  <Button
                    style={{ width: "250px" }}
                    size="field"
                    onClick={onCreateClicked}
                    disabled={shouldCreateBeDisabled() || creating}
                    kind="primary"
                  >
                    Create
                  </Button>
                )}
                <ModalWrapper
                  disabled={shouldSchedulingBeDisabled()}
                  hasForm
                  aria-label="modal"
                  hasScrollingContent
                  buttonTriggerText="Schedule"
                  triggerButtonKind="tertiary"
                  handleSubmit={() => onScheduleSubmit()}
                  shouldCloseAfterSubmit
                >
                  <TextInput.PasswordInput
                    id="api_key_input"
                    labelText="API Key"
                    disabled={apiKeyValid}
                    value={apiKey}
                    placeholder="Enter a valid api key for this account"
                    onChange={(e) =>
                      e.target ? setApiKey(e.target.value) : setApiKey("")
                    }
                  />
                  <Spacer height="16px" />
                  <div>
                    {!apiKeyValid ? (
                      <Button onClick={onSubmitAPIKeyClicked} size="small">
                        Save API Key
                      </Button>
                    ) : (
                      <Button
                        kind="danger"
                        size="small"
                        onClick={onDeleteAPIKeyClicked}
                      >
                        Delete API Key
                      </Button>
                    )}
                  </div>

                  <Spacer height="16px" />
                  <DatePicker
                    dateFormat="m/d/y"
                    datePickerType="range"
                    minDate={setMinDate()}
                    maxDate={setMaxDate()}
                    onChange={(e) => setDateRange(e)}
                  >
                    <DatePickerInput
                      labelText="Start Date"
                      id="start_date_picker"
                      invalid={dateRange.length < 1}
                      invalidText="Need a start date"
                    />
                    <DatePickerInput
                      labelText="End Date"
                      id="end_date_picker"
                      invalid={dateRange.length < 1}
                      invalidText="Need a end date"
                    />
                  </DatePicker>
                  <Spacer height="16px" />
                  <div
                    style={{
                      display: "flex",
                      flexDirection: "row",
                      alignItems: "flex-start",
                    }}
                  >
                    <TimePicker
                      id="start_time_picker"
                      labelText="Start Time"
                      placeholder="hh:mm"
                      type="text"
                      value={startTime}
                      onChange={(e) =>
                        e.target
                          ? setStartTime(e.target.value)
                          : setStartTime("00:00")
                      }
                      invalid={timeInvalid(startTime)}
                      invalidText="invalid time"
                      maxLength={5}
                    >
                      <TimePickerSelect
                        labelText="AM/PM"
                        id="start_time_am_pm"
                        value={startTimeAMPM}
                        onChange={(e) =>
                          e.target ? setStartTimeAMPM(e.target.value) : null
                        }
                      >
                        <SelectItem text="AM" value="AM" />
                        <SelectItem text="PM" value="PM" />
                      </TimePickerSelect>
                    </TimePicker>
                    <TimePicker
                      id="end_time_picker"
                      labelText="End Time"
                      placeholder="hh:mm"
                      type="text"
                      value={endTime}
                      onChange={(e) =>
                        e.target
                          ? setEndTime(e.target.value)
                          : setEndTime("00:00")
                      }
                      invalid={timeInvalid(endTime)}
                      invalidText="invalid time"
                      maxLength={5}
                    >
                      <TimePickerSelect
                        labelText="AM/PM"
                        id="end_time_am_pm"
                        value={endTimeAMPM}
                        onChange={(e) =>
                          e.target ? setEndTimeAMPM(e.target.value) : null
                        }
                      >
                        <SelectItem text="AM" value="AM" />
                        <SelectItem text="PM" value="PM" />
                      </TimePickerSelect>
                    </TimePicker>
                  </div>
                  <Spacer height="16px" />
                  <NotificationEmail
                    accountID={accountID}
                    setSelectedEmails={setSelectedEmails}
                  />
                  <WorkshopAccount
                    isWorkshop={isWorkshop}
                    setIsWorkshop={setIsWorkshop}
                    githubIssue={githubIssue}
                    setGithubIssue={setGithubIssue}
                    userPerCluster={userPerCluster}
                    setUserPerCluster={setUserPerCluster}
                    accountID={accountID}
                  />
                </ModalWrapper>
                <Spacer height="16px" />
              </div>
            </Column>
          </Row>
        </Grid>
      </Form>
    </>
  );
};

export default CreateForm;
