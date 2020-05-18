const headers = [
  {
    key: 'name',
    header: 'Name',
  },
  {
    key: 'state',
    header: 'State',
  },
  {
    key: 'masterKubeVersion',
    header: 'Master Version',
  },
  {
    // `key` is the name of the field on the row object itself for the header
    key: 'location',
    // `header` will be the name you want rendered in the Table Header
    header: 'Location',
  },
  {
    key: 'dataCenter',
    header: 'Data Center',
  },
  {
    key: 'workerCount',
    header: 'Worker Count',
  },
  {
    key: 'tags',
    header: 'Tags',
  },
  {
    key: 'cost',
    header: 'Cost',
  },
];

export default headers;
