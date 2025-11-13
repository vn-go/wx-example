<script lang="ts">
    import Layout from '@layouts/main.svelte';
   
    import Header from '@components/header.svelte';
    let username="abc"
    
    
</script> 
<!-- Gán nội dung cho các slot bằng attribute slot="..." -->
<Layout title="Dashboard">
    <!-- header slot: nhận slot-prop `title` từ layout qua let:title -->
    <svelte:fragment slot="header">
        <!-- <div class="bg-amber-400">
            dsad dsada ddasdsad 
        </div> -->
        <Header/>
    </svelte:fragment>
    
    <svelte:fragment slot="sidebar">
        <div style="min-width:300px">

        </div>
    </svelte:fragment>
    <div></div>
  
  </Layout>