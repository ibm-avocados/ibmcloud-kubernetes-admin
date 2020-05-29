import React from "react";
import { ContentSwitcher, Switch } from "carbon-components-react";

const Calendar = React.lazy(() => import("./Calendar"));

const ScheduleContentSwitcher = ({ accountID }) => {
  const [selected, setSelected] = React.useState(0);

  return (
    <>
      <ContentSwitcher selectedIndex={0} onChange={(i) => setSelected(i.index)}>
        <Switch name="one" text="List" />
        <Switch name="two" text="Calendar" />
      </ContentSwitcher>
      {selected === 0 ? <h1>{accountID}</h1> : null}
      {selected === 1 ? <Calendar accountID={accountID} /> : null}
    </>
  );
};

export default ScheduleContentSwitcher;
