<script>
    import { profileData } from "$lib/state.svelte";
    import { onMount } from "svelte";
    import BottomNav from "../components/BottomNav.svelte";
    import SubscribeForm from "../components/SubscribeForm.svelte";
    import "../global.css";
    import Alert from "../components/Alert.svelte";
    let { children } = $props();

    async function getProfileData() {
        try {
            let response = await fetch("http://localhost:8080/api/user", {
                credentials: "include",
            });
            if (response.status == 302) {
                console.log("redirect");
                window.location.href = await response.text();
                return;
            }
            if (response.status != 200) {
                let err = await response.text();
                throw new Error(`unable to create feed: ${err}`);
            }
            let data = await response.json();
            console.log("data:", data);
            profileData.mashboardEmail = data.email;
            profileData.name = data.firstName + " " + data.lastName;
            profileData.oauthProvider = data.oauthProvider;
        } catch (e) {
            console.error(e);
        }
    }
    onMount(async () => {
        await getProfileData();
    });
</script>

<svelte:head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width" />
    <link rel="icon" href="data:," />
    <link rel="stylesheet" href="https://cdn.plyr.io/3.7.8/plyr.css" />
    <link rel="preconnect" href="https://fonts.googleapis.com" />
    <link rel="preconnect" href="https://fonts.gstatic.com" />
    <link
        href="https://fonts.googleapis.com/css2?family=Rubik:ital,wght@0,300..900;1,300..900&display=swap"
        rel="stylesheet"
    />
    <script
        src="https://kit.fontawesome.com/1779bb385a.js"
        crossorigin="anonymous"
    ></script>
    <title>Mashboard</title>
</svelte:head>

{@render children()}

<div id="alert-container" class="absolute bottom-20 left-0 w-full z-20">
    <Alert />
</div>
<SubscribeForm />
<BottomNav />
