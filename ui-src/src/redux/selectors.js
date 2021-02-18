
export const getApps = state => state.apps.data;
export const getAppsPending = state => state.apps.pending;
export const getAppsError = state => state.apps.error;

export const getUsers = state => state.users.data;
export const getUsersPending = state => state.users.pending;
export const getUsersError = state => state.users.error;

export const getStatus = state => state.status.data;
export const getStatusPending = state => state.status.pending;
export const getStatusError = state => state.status.error;

export const getSidebarShow = state => state.ui.sidebarShow;

export const getAdminApp = state => state.apps.data.find( app => app.id === 1 );

