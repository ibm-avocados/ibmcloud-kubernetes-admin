import React from 'react';
import { TextInput, Button, MultiSelect } from 'carbon-components-react';

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
  const [emailText, setEmailText] = React.useState('');
  const [updateEmailText, setUpdateEmailText] = React.useState('');

  React.useEffect(() => {
    getAccountAdmins();
  }, [accountID]);

  const getAccountAdmins = async () => {
    try {
      const emails = await grab(`/api/v1/notification/${accountID}/email`);
      setEmailAvailable(true);
      setEmails(emails);
    } catch (e) {
      console.log(e);
    }
  };

  const onEmailSubmit = async () => {
    const emails = emailText
      .toLowerCase()
      .split(',')
      .map((email) => email.trim());
    try {
      const response = await grab('/api/v1/notification/email/create', {
        method: 'post',
        body: JSON.stringify({
          email: emails,
          accountID: accountID,
        }),
      });
      console.log('email submitted');
      setEmails(emails);
      setEmailAvailable(true);
    } catch (e) {
      console.log(e);
    }
  };

  const onAddEmailSubmit = async () => {
    const emails = updateEmailText
      .toLowerCase()
      .split(',')
      .map((email) => email.trim());
    try {
      const response = await grab('/api/v1/notification/email/add', {
        method: 'put',
        body: JSON.stringify({
          email: emails,
          accountID: accountID,
        }),
      });
      console.log('email added');
      getAccountAdmins();
      setUpdateEmailText('');
    } catch (e) {
      console.log(e);
    }
  };

  const onRemoveEmailSubmit = async () => {
    const emails = updateEmailText
      .toLowerCase()
      .split(',')
      .map((email) => email.trim());
    try {
      const response = await grab('/api/v1/notification/email/remove', {
        method: 'put',
        body: JSON.stringify({
          email: emails,
          accountID: accountID,
        }),
      });
      console.log('email removed');
      getAccountAdmins();
      setUpdateEmailText('');
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

  const onEmailSelected = (selectedItems) => {
    if (!selectedItems || selectedItems.length === 0) {
      setSelectedEmails([]);
    }
    setSelectedEmails(selectedItems.map((v) => v.label));
  };

  return (
    <>
      {!emailAvailable ? (
        <>
          <TextInput
            id="email-input"
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
            id="email-select"
            titleText="Notify emails"
            items={emailItems}
            itemToString={(item) => item.text}
            onChange={({ selectedItems }) => onEmailSelected(selectedItems)}
            label="Select specific email addresses to notify (defaults to all)"
          />
          <Spacer height="16px" />
          <TextInput
            id="add-email-input"
            value={updateEmailText}
            onChange={(e) => setUpdateEmailText(e.target.value.trim())}
            labelText="Update admin emails"
            placeholder="user1@email.com,user2@email.com"
          />
          <Spacer height="16px" />
          <Button
            onClick={onAddEmailSubmit}
            size="small"
            disabled={updateEmailText.trim().length < 6}
          >
            Add Email for account
          </Button>
          <Button
            onClick={onRemoveEmailSubmit}
            size="small"
            kind="danger"
            disabled={updateEmailText.trim().length < 6}
          >
            Remove Email for account
          </Button>
        </>
      )}

      <Spacer height="16px" />
    </>
  );
};

export default NotificationEmail;
