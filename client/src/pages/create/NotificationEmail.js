import React from "react";
import { TextInput, Button, MultiSelect } from "carbon-components-react";

const Spacer = ({ height }) => <div style={{ marginTop: height }} />;

const grab = async (url, options, retryCount = 0) => {
  const response = await fetch(url, options);
  const data = await response.json();
  if (response.status !== 200) {
    if (retryCount > 0) {
      return await grab(url, options, retryCount - 1);
    }
    throw Error(data);
  }

  return data;
};

const NotificationEmail = ({ accountID, setSelectedEmails }) => {
  const [emailAvailable, setEmailAvailable] = React.useState(false);
  const [emails, setEmails] = React.useState([]);
  const [emailText, setEmailText] = React.useState("");

  React.useEffect(() => {
    const getAccountAdmins = async () => {
      try {
        const emails = await grab(`/api/v1/notification/${accountID}/email`);
        setEmailAvailable(true);
        setEmails(emails);
      } catch (e) {
        console.log(e);
      }
    };
    getAccountAdmins();
  }, [accountID]);

  const onEmailSubmit = async () => {
    const emails = emailText
      .toLowerCase()
      .split(",")
      .map((email) => email.trim());
    try {
      const response = await grab("/api/v1/notification/email/create", {
        method: "post",
        body: JSON.stringify({
          email: emails,
          accountID: accountID,
        }),
      });
      console.log(response);
      setEmails(emails);
      setEmailAvailable(true);
    } catch (e) {
      console.log(e);
    }
  };

  const emailItems = emails
    ? emails.map((email, i) => ({
        id: i,
        text: email,
        label: email,
      }))
    : [];

  return (
    <>
      {!emailAvailable ? (
        <>
          <TextInput
            value={emailText}
            onChange={(e) => setEmailText(e.target.value.trim())}
            labelText="Account admin emails"
            placeholder="user1@email.com,user2@email.com"
          />
          <Spacer height="16px" />
          <Button
            onClick={onEmailSubmit}
            size="small"
            disabled={emailText.trim().length < 6}
          >
            Save Email for account
          </Button>
        </>
      ) : (
        <>
          <MultiSelect
            titleText="Notify emails"
            items={emailItems}
            itemToString={(item) => item.text}
            onChange={(items) => setSelectedEmails(items.map(item => item.text))}
            label="Select email addresses"
          />
        </>
      )}

      <Spacer height="16px" />
    </>
  );
};

export default NotificationEmail;
