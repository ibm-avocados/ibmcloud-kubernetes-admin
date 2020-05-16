import React from "react";

import {
  Form,
  TextInput,
  Button,
  Dropdown,
  RadioTile,
  Row,
  Grid,
  Column,
  FormLabel,
  Tooltip,
  InlineLoading,
  ToastNotification,
  ModalWrapper,
  DatePicker,
  DatePickerInput,
  TimePicker,
  TimePickerSelect,
  SelectItem,
} from "carbon-components-react";

import geos from "../data/geo";

import styles from "./CreateForm.module.css";

import "./CreateForm.css";

const Spacer = ({ height }) => <div style={{ marginTop: height }} />;

const Divider = ({ width }) => <div style={{ marginRight: width }} />;

const grab = async (url, options, retryCount = 0) => {
  const response = await fetch(url, options);
  const data = await response.json();
  if (response.status !== 200) {
    if (retryCount > 0) {
      return await grab(url, options, retryCount - 1);
    }
    throw Error(data);
  }

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
  const [privateVlans, setPrivateVlans] = React.useState([]);
  const [publicVlans, setPublicVlans] = React.useState([]);
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
  const [selectedWorkerZone, setSelectedWorkerZone] = React.useState(null);
  const [selectedPrivateVlan, setSelecetedPrivateVlan] = React.useState(null);
  const [selectedPublicVlan, setSelecetedPublicVlan] = React.useState(null);
  const [selectedFlavor, setSelectedFlavor] = React.useState(null);
  const [selectedGroup, setSelectedGroup] = React.useState(null);
  // ui indicators
  const [creating, setCreating] = React.useState(false);
  const [loaderDescription, setLoaderDescription] = React.useState("");
  const [createSuccess, setCreateSuccess] = React.useState(false);
  const [startTimeAMPM, setStartTimeAMPM] = React.useState("AM");
  const [endTimeAMPM, setEndTimeAMPM] = React.useState("AM");
  const [apiKey, setApiKey] = React.useState("");
  const [apiKeyValid, setApiKeyValid] = React.useState(false);
  const [startTime, setStartTime] = React.useState("");
  const [endTime, setEndTime] = React.useState("");
  const [dateRange, setDateRange] = React.useState([]);

  React.useEffect(() => {
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
        console.log(resourceGroups);
        if (resourceGroups) {
          setResourceGroups(resourceGroups.resources);
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
  }, [accountID]);

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

  const getVlans = async (datacenter) => {
    try {
      const vlans = await grab(`/api/v1/clusters/${datacenter}/vlans`);
      const privateVlans = vlans.filter((vlan) => vlan.type === "private");
      if (privateVlans && privateVlans.length > 0) {
        setPrivateVlans(privateVlans);
        setSelecetedPrivateVlan(privateVlans[0]);
      }
      const publicVlans = vlans.filter((vlan) => vlan.type === "public");
      if (publicVlans && publicVlans.length > 0) {
        setPublicVlans(publicVlans);
        setSelecetedPublicVlan(publicVlans[0]);
      }
    } catch (e) {
      console.log(e);
    }
  };

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
    setSelectedWorkerZone(null);
    setPrivateVlans([]);
    setSelecetedPrivateVlan(null);
    setPublicVlans([]);
    setSelecetedPublicVlan(null);
  };

  const onWorkerZoneSelected = (zone) => {
    setSelectedWorkerZone(zone);
    setPrivateVlans([]);
    setSelecetedPrivateVlan(null);
    setPublicVlans([]);
    setSelecetedPublicVlan(null);
    getVlans(zone.id);
    getFlavors(zone.id);
  };

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
        dataCenter: selectedWorkerZone.id,
        defaultWorkerPoolName: "",
        defaultWorkerPoolEntitlement,
        disableAutoUpdate: true,
        noSubnnet: false,
        podSubnet: "",
        serviceSubnet: "",
        machineType: selectedFlavor.name,
        privateVlan: selectedPrivateVlan.id,
        publicVlan: selectedPublicVlan.id,
        masterVersion: version,
        workerNum: Number(workerCount),
        diskEncryption: true,
        isolation: "public",
        GatewayEnabled: false,
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

    for (let i = 0; i < range; i++) {
      setLoaderDescription(`Creating Cluster ${i + 1} of ${range + 1}`);
      console.log("Creating luster ", i);

      const CreateClusterRequest = request[i];

      try {
        const clusterResponse = await grab("/api/v1/clusters", {
          method: "post",
          body: JSON.stringify(CreateClusterRequest),
        });

        console.log(clusterResponse);

        console.log("Sleeping 5s before trying to set tags");
        setLoaderDescription(
          `Preparing to Tag Cluster ${i + 1} of ${range + 1}`
        );
        await sleep(5000);
        setLoaderDescription(`Tagging Cluster ${i + 1} of ${range + 1}`);

        // comma separated tags.
        const tagPromises = tags.split(",").map(async (tag) => {
          try {
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
        console.log(e);
      }
    }

    setCreateSuccess(true);
    setCreating(false);
    // console.log(JSON.stringify(CreateClusterRequest));
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
    const zoneSelected = !!selectedWorkerZone;
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
    return false;
  };

  const shouldScheduleSubmitBeDisabled = () => {
    const dateSet = dateRange.length === 2;
    const createTimetSet = !timeInvalid(startTime);
    const endTimeSet = !timeInvalid(endTime);

    return !(dateSet && createTimetSet && endTimeSet && apiKeyValid);
  };

  Date.prototype.addTime = function(h, m) {
    this.setTime(this.getTime() + (h*60*60*1000) + m * 60 * 1000);
    return this;
  }

  const onScheduleSubmit = () => {
    // const request = getCreateRequest();
    let startDate = dateRange[0];
    let endDate = dateRange[1];
    let startHour = Number(startTime.split(":")[0])
    console.log(startHour);
    const startMinute = Number(startTime.split(":")[1])
    const offset = startDate.getTimezoneOffset();
    console.log(startMinute);
    startHour += startTimeAMPM === "PM"? 12 : 0;
    console.log(startHour);
    let endHour = Number(endTime.split(":")[0])
    const endMinute = Number(endTime.split(":")[1])
    endHour += startTimeAMPM === "PM"? 12 : 0;

    startDate.addTime(startHour, startMinute - offset);
    const createAt = startDate.getTime()/1000;

    

    console.log(createAt);
    console.log(startDate);
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
          </Row>
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
              <FormLabel>Cluster count</FormLabel>
              <TextInput
                value={clusterCount}
                onChange={(e) => setClusterCount(e.target.value.trim())}
                labelText=""
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

          {createSuccess ? (
            <>
              <Spacer height="16px" />
              <ToastNotification
                title="Cluster Created"
                subtitle={`${clusterCount} ${
                  kubernetesSelected ? "Kubernetes" : "Openshift"
                } Clusters Created`}
                kind="success"
                caption={`Datacenter: ${selectedWorkerZone.display_name}, ${selectedRegion.display_name}`}
                timeout={5000}
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
                        onChange={(e) => e.target?setStartTimeAMPM(e.target.value):null}
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
                        onChange={(e) => e.target?setEndTimeAMPM(e.target.value):null}
                      >
                        <SelectItem text="AM" value="AM" />
                        <SelectItem text="PM" value="PM" />
                      </TimePickerSelect>
                    </TimePicker>
                  </div>
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
