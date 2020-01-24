import React from "react";
import Login from "./Login";
import TextInput from "./common/TextInput";
import Navbar from "./common/Navbar";
import Clusters from "./common/Clusters";

class Home extends React.Component {
  constructor(props) {
    super(props);
    this.state = { isLoggedIn: 0 };
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

    this.setState({items: accounts.resources})
  }

  accountSelected = async ({selectedItem}) => {
    console.log(selectedItem);
    let response = await fetch("/api/v1/authenticate/account", {
      method: "POST",
      body: JSON.stringify({
        id: selectedItem.metadata.guid
      })
    });
    if (response.status === 200) {
      this.loadClusters()
    }
  }

  loadClusters = async () => {
    let clusterResponse = await fetch("/api/v1/clusters");
    let clusters = await clusterResponse.json()
    this.setState({clusters: clusters});
    console.log(this.state.clusters);
  }


  render() {
    if (this.state.isLoggedIn === 2) {
      return (
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
    }

    if (this.state.isLoggedIn === 1) {
      return (
        <>
          <Navbar items={this.state.items} accountSelected={this.accountSelected}/>
          {this.state.clusters ? <Clusters data={this.state.clusters}/> : <></>}
        </>
      );
    }
    return (
      <>
        <Login onResult={this.handleLoginResponse} getAccounts={this.loadAccounts}/>
      </>
    );
  }
}

export default Home;
