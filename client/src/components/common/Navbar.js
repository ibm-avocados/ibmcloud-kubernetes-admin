import React from "react";
import styles from "./Navbar.module.css";
import { Dropdown } from "carbon-components-react";
import "./Dropdown.css";

const Navbar = props => {
  const itemToString = item => {
    if (item) {
      let name = item.entity.name;
      let softlayerAccountId =
        item.entity.bluemix_subscriptions[0].softlayer_account_id || "";

      return name + " " + softlayerAccountId;
    }
    return "Unknown";
  };


  return (
    <>
      <div className={styles.wrapper}>
        <div className={styles.title}>
          <span className={styles.bold}>IBM</span> Cloud
        </div>

        <Dropdown
          disabled={props.accountsLoaded}
          className={styles.dropdown}
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
