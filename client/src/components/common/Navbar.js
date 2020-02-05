import React from "react";
import styles from "./Navbar.module.css";
import { Dropdown } from "carbon-components-react";
import history from "../../globalHistory";
import "./Dropdown.css";

const MenuItem = props => {
  return (
    <>
      <div
        className={styles.menuItem}
        style={props.style}
        onClick={props.onClickHandler}
      >
        {props.label}
      </div>
    </>
  );
};

const Navbar = ({ items, accountsLoaded, accountSelected }) => {
  const itemToString = item => {
    if (item) {
      let name = item.entity.name;
      let softlayerAccountId =
        item.entity.bluemix_subscriptions[0].softlayer_account_id || "";

      return name + " " + softlayerAccountId;
    }
    return "Unknown";
  };

  const handleClick = () => {
    history.push("/create");
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
        <MenuItem label="Create" onClickHandler={handleClick} />
        <Dropdown
          disabled={accountsLoaded}
          className={styles.dropdown}
          ariaLabel="Dropdown"
          label="Select Account"
          items={items || []}
          onChange={accountSelected}
          itemToString={itemToString}
          id="account-dropdown"
          light={false}
        />
      </div>
    </>
  );
};

export default Navbar;
