import React from "react";
import Clusters from "./common/Clusters";
import {Loading} from 'carbon-components-react';

const AppPage = ({ hasChosenAccount, tokenUpgraded, accountID }) => {
  if (!hasChosenAccount) {
    return null;
  } else if (tokenUpgraded) {
    return <Clusters accountID={accountID} />;
  } else {
    return <Loading />
  }
};
export default AppPage;
