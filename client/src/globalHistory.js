import { createBrowserHistory } from 'history';

const forceHistory = createBrowserHistory({forceRefresh:true})

export default createBrowserHistory();
export {forceHistory}
