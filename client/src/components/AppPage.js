import React from "react";
import Clusters from "./common/Clusters";

const AppPage = ({ hasChosenAccount, tokenUpgraded, accountID }) => {
  return (
    <>
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
  } else if (tokenUpgraded) {
    return <Clusters accountID={accountID} />;
  } else {
    return <h1>Token Not Valid</h1>
  }
};
export default AppPage;
