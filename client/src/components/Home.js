import React from "react";
import Login from "./Login";
import TextInput from "./common/TextInput";
import Navbar from "./common/Navbar";
import Clusters from "./common/Clusters";
import { DataTableSkeleton } from "carbon-components-react";

class Home extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      isLoggedIn: 0,
      clusterLoaded: false,
      accountsLoaded: false,
      accountSelected: false
    };
  }
  handleLoginResponse = loggedIn => {
    if (loggedIn !== this.state.isLoggedIn) {
      this.setState({ isLoggedIn: loggedIn });
    }
  };

  handleOtp = event => {
    this.setState({
      otp: event.target.value
    });
  };

  onKeyDown = async event => {
    if (event.key === "Enter") {
      event.preventDefault();
      event.stopPropagation();
      console.log("Enter Pressed");
      let response = await fetch("/api/v1/authenticate", {
        method: "POST",
        body: JSON.stringify({
          otp: this.state.otp
        })
      });
      if (response.status > 400) {
        this.setState({ isLoggedIn: 0 });
      } else if (response.status === 200) {
        this.loadAccounts();
        this.setState({ isLoggedIn: 1 });
      }
    }
  };

  loadAccounts = async () => {
    let response = await fetch("/api/v1/accounts");
    let accounts = await response.json();

    this.setState({ accountsLoaded: true });
    this.setState({ accounts: accounts.resources });
  };

  accountSelected = async ({ selectedItem }) => {
    console.log(selectedItem);
    let response = await fetch("/api/v1/authenticate/account", {
      method: "POST",
      body: JSON.stringify({
        id: selectedItem.metadata.guid
      })
    });
    if (response.status === 200) {
      this.setState({ accountSelected: true });
      this.loadClusters();
    }
  };

  loadClusters = async () => {
    this.setState({ clusterLoaded: false });
    let clusterResponse = await fetch("/api/v1/clusters");
    let clusters = await clusterResponse.json();
    this.setState({ clusters: clusters });
    this.setState({ clusterLoaded: true });
  };

  otp = () => (
    <>
      <TextInput
        style={{
          position: "absolute",
          left: "50%",
          top: "50%",
          transform: "translate(-50%, -50%)"
        }}
        onChange={this.handleOtp}
        placeholder="One Time Passcode"
        onKeyDown={this.onKeyDown}
      />
    </>
  );

  showCluster = () => (
    <>
      {this.state.clusterLoaded ? (
        <Clusters data={this.state.clusters} />
      ) : (
        <DataTableSkeleton
          columnCount={6}
          compact={false}
          headers={[
            "Name",
            "State",
            "Master Version",
            "Location",
            "Data Center",
            "Worker Count"
          ]}
          rowCount={5}
          zebra={false}
        />
      )}
    </>
  );

  loggedIn = () => {
    return (
      <>
        <Navbar
          isLoaded={this.state.accountsLoaded}
          items={this.state.accounts}
          accountSelected={this.accountSelected}
        />
        {this.state.accountSelected ? this.showCluster() : <></>}
      </>
    );
  };

  render() {
    if (this.state.isLoggedIn === 2) {
      return this.otp();
    }

    if (this.state.isLoggedIn === 1) {
      let data = this.loggedIn();
      return data;
    }

    return (
      <>
        <Login
          onResult={this.handleLoginResponse}
          getAccounts={this.loadAccounts}
        />
      </>
    );
  }
}

export default Home;
