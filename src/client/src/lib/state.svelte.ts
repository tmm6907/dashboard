
export const feedState = $state({
    feedItems: [],
    feedCollections: [],
    feedLatest: [],
    feedSaved: [],
    category: "all"
});

export const selectedItem = $state({
    id: 1,
    data: {}
})


export const profileData = $state({
    name: "",
    mashboardEmail: "",
    oauthProvider: "",
})


