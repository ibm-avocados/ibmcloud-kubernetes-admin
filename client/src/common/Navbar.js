import React from "react";
import { Dropdown } from "carbon-components-react";
import styles from "./Navbar.module.css";
// import Dropdown from "react-dropdown";
import history from "../globalHistory";
import "./Dropdown.css";

const MenuItem = (props) => {
  return (
    <>
      <div
        className={styles.menuItem}
        onClick={props.onClickHandler}
      >
        {props.label}
      </div>
    </>
  );
};

const Navbar = (props) => {
  const itemToString = (item) => {
    if (item) {
      const { name } = item.entity;
      const softlayerAccountId =
        item.entity.bluemix_subscriptions[0].softlayer_account_id || "";

      return `${name} ${softlayerAccountId}`;
    }
    return "Unknown";
  };

  const handleCreateClick = () => {
    history.push("/create");
  };
  const handleScheduleClick = () => {
    history.push("/schedule");
  };

  const homeClick = () => {
    history.push("/");
  };

  return (
    <>
      <div className={styles.wrapper}>
        <div className={styles.title} onClick={homeClick}>
          <span className={styles.bold}>IBM</span> Cloud
        </div>
        <MenuItem label="Create" onClickHandler={handleCreateClick} />
        <MenuItem label="Schedule" onClickHandler={handleScheduleClick} />
        <Dropdown
          disabled={props.accountsLoaded}
          className="navbar-dropdown"
          ariaLabel="Dropdown"
          label="Select Account"
          items={props.items || []}
          onChange={props.accountSelected}
          itemToString={itemToString}
          id="account-dropdown"
          light={false}
        />
      </div>
    </>
  );
};

export default Navbar;
