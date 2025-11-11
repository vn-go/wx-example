<script lang="ts">
    import Layout from '@layouts/main.svelte'
    let username="abc"
</script> 
<!-- Gán nội dung cho các slot bằng attribute slot="..." -->
<Layout title="Dashboard">
    <!-- header slot: nhận slot-prop `title` từ layout qua let:title -->
    <svelte:fragment slot="header">
        <div class="bg-amber-400">
            dsad dsada ddasdsad 
        </div>
        
    </svelte:fragment>
    
    <svelte:fragment slot="sidebar">
        <div style="min-width:300px">

        </div>
    </svelte:fragment>
    <div></div>
  
  </Layout>