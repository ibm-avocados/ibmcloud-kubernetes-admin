import React, { useCallback, useEffect, useState } from "react";
import Navbar from "./common/Navbar";
import Clusters from "./common/Clusters";
import { getJSON } from "../fetchUtil";


const AppPage = () => {
  const [isLoadingAccounts, setLoadingAccounts] = useState(true);
  const [accounts, setAccounts] = useState([]);
  const [accountID, setSelectedAccountID] = useState();


  const [hasChosenAccount, setHasChosenAccount] = useState(false);
  const [tokenUpgraded, setTokenUpgraded] = useState(false);

  const loadAccounts = async () => {
    setLoadingAccounts(true);
    const accounts = await fetch("/api/v1/accounts").then(getJSON);
    setAccounts(accounts.resources);
    setLoadingAccounts(false);
  };

  useEffect(() => {
    // You can't have an async function as an effect argument.
    loadAccounts();
  }, []);

  const handleAccountChosen = useCallback(
    async ({ selectedItem }) => {
      setSelectedAccountID(selectedItem.metadata.guid);
      setTokenUpgraded(false);
      setHasChosenAccount(true);
      const { status } = await fetch("/api/v1/authenticate/account", {
        method: "POST",
        body: JSON.stringify({
          id: selectedItem.metadata.guid
        })
      });
      if (status === 200) {
        setTokenUpgraded(true);
      }
    },
    []
  );

  return (
    <>
      <Navbar
        isLoaded={!isLoadingAccounts}
        items={accounts}
        accountSelected={handleAccountChosen}
      />
      <ConditionalClusterTable
        accountChanged={hasChosenAccount}
        tokenUpgraded={tokenUpgraded}
        accountID={accountID}
      />
    </>
  );
};

const ConditionalClusterTable = ({ accountChanged, tokenUpgraded, accountID }) => {
  if (!accountChanged) {
    return null;
  } else if(tokenUpgraded){
    return <Clusters accountID={accountID}/>;
  } else {
    return <h1>Token Not Valid</h1>
  }

  
};
export default AppPage;
