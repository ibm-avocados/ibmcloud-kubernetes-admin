import React from "react";
import { Loading } from "carbon-components-react";
import ScheduleContentSwitcher from './ScheduleContentSwitcher';

const SchedulePage = ({ hasChosenAccount, tokenUpgraded, accountID }) => {
  if (!hasChosenAccount) {
    return <h1>Please select account</h1>;
  } else if (tokenUpgraded) {
    return <ScheduleContentSwitcher />;
  }
  return <Loading />;
};

export default SchedulePage;
