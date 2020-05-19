import React from "react";
import FullCalendar from "@fullcalendar/react";
import dayGridPlugin from "@fullcalendar/daygrid";
import timeGridPlugin from "@fullcalendar/timegrid";

import "./main.scss"; // webpack must be configured to do this

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

function getWindowDimensions() {
  const { innerWidth: width, innerHeight: height } = window;
  return {
    width,
    height,
  };
}

const ScheduleCalendar = ({ accountID }) => {
  const [windowDimensions, setWindowDimensions] = React.useState(
    getWindowDimensions()
  );

  const [events, setEvents] = React.useState([]);

  React.useEffect(() => {
    const handleResize = () => {
      setWindowDimensions(getWindowDimensions());
    };

    const loadEvents = async () => {
      try {
        const data = await grab(`/api/v1/schedule/${accountID}/all`);
        console.log(data);
        const events = data.map((schedule, i) => {
          const start = new Date(schedule["createAt"] * 1000).toISOString();
          const end = new Date(schedule["destroyAt"] * 1000).toISOString();
          let _title = schedule["ClusterRequests"][0]["clusterRequest"]["name"];
          const title = _title.substring(0, _title.length-4);


          const id = schedule["_id"];

          const status = schedule["status"];
          let bg = "#0f62fe";
          switch(status) {
            case "scheduled":
              bg = "#0f62fe";
              break;
            case "created":
              bg = "#24a148";
              break;
            case "completed":
              bg = "#525252";
              break;
            default:
              bg = "#0f62fe";
              break;
          }
          return {
            id: id,
            title: title,
            start: start,
            end: end,
            backgroundColor: bg,
          }
        })
        console.log(events);
        setEvents(events);
      } catch (e) {
        console.log(e);
      }
    };

    loadEvents();

    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, [accountID]);
  return (
    <FullCalendar
      schedulerLicenseKey="GPL-My-Project-Is-Open-Source"
      height={windowDimensions.height - 100}
      header={{
        left: "prev,next today",
        center: "title",
        right: "dayGridMonth,timeGridWeek",
      }}
      events={events}
      defaultView="dayGridMonth"
      plugins={[dayGridPlugin, timeGridPlugin]}
    />
  );
};

export default ScheduleCalendar;
