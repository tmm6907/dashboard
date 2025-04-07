<script>
    import { onMount } from "svelte";
    import NewFeedForm from "../../components/NewFeedForm.svelte";
    function setDataTheme() {
        // Set the data-theme attribute on the <html> element
        let theme = JSON.parse(localStorage.getItem("isToggled") || "false");
    }
    onMount(() => {
        setDataTheme();
    });
</script>

<div class="p-8">
    <header aria-labelledby="new-feed-heading">
        <h1 id="new-feed-heading" class="text-xl text-secondary">
            Create New Feed
        </h1>
    </header>
    <NewFeedForm />
</div>
