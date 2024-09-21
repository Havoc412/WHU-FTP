export type fileItemType = {    // used by MainView.vue && filesPenal.vue
    // path: string;
    name: string;
    size: number;
}

export type WSprepareType = {  // used by MainView.vue && terminalPenal.vue
    type: boolean;
    sessionId: string;

    path: string;
    fileName: string;
    size: number;
};