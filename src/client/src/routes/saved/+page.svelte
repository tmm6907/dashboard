<script>
    import { onMount } from "svelte";
    import Section from "../../components/Section.svelte";
    export const prerender = false;

    let items = $state([]);

    const getSavedItems = async () => {
        try {
            let response = await fetch(
                `http://localhost:8080/api/feeds/items/saved/`,
                {
                    credentials: "include",
                },
            );
            if (response.status == 302) {
                window.location.href = await response.text();
                return;
            }
            if (response.status != 200) {
                throw new Error(await response.text());
            }
            ({ items } = await response.json());
            console.log("saved items:", items);
        } catch (e) {
            console.error(e);
        }
    };

    onMount(async () => {
        console.log("loaded");
        await getSavedItems();
    });
</script>

<div class="p-4">
    <Section heading="Saved" {items} vertical={true} />
</div>
