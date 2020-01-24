import React from "react";
import { DataTable, Button } from "carbon-components-react";
import { Delete16 as Delete, Save16 as Save } from "@carbon/icons-react";

import "./Cluster.css";

const {
  TableContainer,
  Table,
  TableHead,
  TableRow,
  TableBody,
  TableCell,
  TableHeader,
  TableSelectRow,
  TableSelectAll,
  TableToolbar,
  TableToolbarSearch,
  TableToolbarContent,
  TableBatchActions,
  TableBatchAction
} = DataTable;

// We would have a headers array like the following
const headers = [
  {
    key: "name",
    header: "Name"
  },
  {
    key: "state",
    header: "State"
  },
  {
    key: "masterKubeVersion",
    header: "Master Version"
  },
  {
    // `key` is the name of the field on the row object itself for the header
    key: "location",
    // `header` will be the name you want rendered in the Table Header
    header: "Location"
  },
  {
    key: "dataCenter",
    header: "Data Center"
  },
  {
    key: "workerCount",
    header: "Worker Count"
  }, 
  {
    key: "crn",
    header: "Crn"
  }
];

const buttonClicked = rows => {
  console.log("slected rows", rows);
};

const render = ({
  rows,
  headers,
  getHeaderProps,
  getBatchActionProps,
  getSelectionProps,
  filterRows,
  selectedRows,
  onInputChange
}) => {
  return (
    <TableContainer title="Clusters">
      <TableToolbar>
        {/* pass in `onInputChange` change here to make filtering work */}
        <TableBatchActions {...getBatchActionProps()}>
          <TableBatchAction
            tabIndex={getBatchActionProps().shouldShowBatchActions ? 0 : -1}
            renderIcon={Delete}
            onClick={buttonClicked(() => selectedRows)}
          >
            Delete
          </TableBatchAction>
          <TableBatchAction
            renderIcon={Save}
            onClick={() => buttonClicked(selectedRows)}
          >
            Save
          </TableBatchAction>
        </TableBatchActions>
        <TableToolbarContent>
          <TableToolbarSearch
            tabIndex={getBatchActionProps().shouldShowBatchActions ? -1 : 0}
            onChange={onInputChange}
          />
        </TableToolbarContent>
        {/* <TableToolbarContent>
          <Button onClick={() => buttonClicked(selectedRows)}  kind="primary">
            Add new
          </Button>
        </TableToolbarContent> */}
      </TableToolbar>
      <Table>
        <TableHead>
          <TableRow>
            <TableSelectAll {...getSelectionProps()} />
            {headers.map(header => (
              <TableHeader {...getHeaderProps({ header })}>
                {processHeader(header)}
              </TableHeader>
            ))}
          </TableRow>
        </TableHead>
        <TableBody>
          {rows.map(row => {
            return (
              <TableRow key={row.id}>
                <TableSelectRow {...getSelectionProps({ row })} />
                {row.cells.map(cell => (
                  <TableCell key={cell.id}>{process(cell)}</TableCell>
                ))}
              </TableRow>
            );
          })}
        </TableBody>
      </Table>
    </TableContainer>
  );
};

const processHeader = header => {
  if(header.key == "crn") {
    return <></>;
  }
  return header.header;
}

const process = cell => {
  let id = cell.id;
  let field = id.split(":")[1];
  let value = cell.value;
  if (field === "state") {
    if (value === "normal") {
      return (
        <>
          <span className="oneline">
            <span className="status normal"></span>
            {value}
          </span>
        </>
      );
    } else if (value === "warning") {
      return (
        <>
          <span className="oneline">
            <span className="status warning"></span>
            {value}
          </span>
        </>
      );
    }
  } else if (field === "masterKubeVersion") {
    if (value.includes("openshift")) {
      return (
        <>
          <span className="oneline">
            <img
              className="logo-image"
              src="https://cloud.ibm.com/kubernetes/img/openshift_logo-7825001afb.svg"
            />
            {value}
          </span>
        </>
      );
    }
    return (
      <>
        <span className="oneline">
          <img
            className="logo-image"
            src="https://cloud.ibm.com/kubernetes/img/container-service-logo-7e87826329.svg"
          />
          {value}
        </span>
      </>
    );
  } else if (field === "crn") {
    return <></>;
  }
  return <>{value}</>;
};

const Clusters = props => {
  return (
    <DataTable
      rows={props.data}
      headers={headers}
      render={render}
      isSortable={true}
      stickyHeader={true}
    />
  );
};

export default Clusters;
