import React from "react";
import { ContentSwitcher, Switch } from "carbon-components-react";

const ScheduleContentSwitcher = () => {
  const [selected, setSelected] = React.useState(0);

  return (
    <>
      <ContentSwitcher selectedIndex={0} onChange={(i) => setSelected(i.index)}>
        <Switch name="one" text="List" />
        <Switch name="two" text="Calendar" />
      </ContentSwitcher>
      {selected === 0 ? <h1>List Showing</h1> : null}
      {selected === 1 ? <h1>Calendar Showing</h1> : null}
    </>
  );
};

export default ScheduleContentSwitcher;
