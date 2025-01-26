<script lang="ts">
    import * as Sidebar from "$lib/components/ui/sidebar/index.js";
    import AppSidebar from "@/components/app_sidebar.svelte";
    
    let { data, children } = $props();
</script>

<Sidebar.Provider>
    <AppSidebar />
    <main class="w-full">
      <div class="absolute bottom-0 left-0 z-50">
        <Sidebar.Trigger />
      </div>
      {@render children?.()}
    </main>
</Sidebar.Provider>