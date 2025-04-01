<script>
    import BottomNav from "../components/BottomNav.svelte";
    import SubscribeForm from "../components/SubscribeForm.svelte";
    import "../global.css";
    let { children } = $props();
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

<SubscribeForm />
<BottomNav />
