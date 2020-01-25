import React, { useCallback, useEffect, useState } from "react";
import Navbar from "./common/Navbar";
import Clusters from "./common/Clusters";
import { DataTableSkeleton } from "carbon-components-react";
import headers from "./data/headers";
import { getJSON } from "../fetchUtil";

const AppPage = () => {
  const [isLoadingAccounts, setLoadingAccounts] = useState(true);
  const [accounts, setAccounts] = useState([]);

  const [hasChosenAccount, setHasChosenAccount] = useState(false);
  const [isLoadingClusters, setLoadingClusters] = useState(true);
  const [clusters, setClusters] = useState([]);

  useEffect(() => {
    // You can't have an async function as an effect argument.
    const loadAccounts = async () => {
      setLoadingAccounts(true);
      const accounts = await fetch("/api/v1/accounts").then(getJSON);
      setAccounts(accounts.resources);
      setLoadingAccounts(false);
    };

    loadAccounts();
  }, []);

  const loadClusters = useCallback(async () => {
    setLoadingClusters(true);
    const clusters = await fetch("/api/v1/clusters").then(getJSON);
    setClusters(clusters);
    setLoadingClusters(false);
  }, []);

  const handleAccountChosen = useCallback(
    async ({ selectedItem }) => {
      setHasChosenAccount(true);
      const { status } = await fetch("/api/v1/authenticate/account", {
        method: "POST",
        body: JSON.stringify({
          id: selectedItem.metadata.guid
        })
      });
      if (status === 200) {
        loadClusters();
      }
    },
    [loadClusters]
  );

  return (
    <>
      <Navbar
        isLoaded={!isLoadingAccounts}
        items={accounts}
        accountSelected={handleAccountChosen}
      />
      <ConditionalClusterTable
        clusters={clusters}
        loading={isLoadingClusters}
        hidden={!hasChosenAccount}
      />
    </>
  );
};

const ConditionalClusterTable = ({ clusters, hidden, loading }) => {
  if (hidden) {
    return null;
  }

  if (loading) {
    return (
      <>
        <div className="bx--data-table-header">
          <h4>Clusters</h4>
        </div>
        <DataTableSkeleton
          columnCount={headers.length}
          compact={false}
          headers={headers}
          rowCount={5}
          zebra={true}
        />
      </>
    );
  }

  return <Clusters headers={headers} data={clusters} />;
};
export default AppPage;
