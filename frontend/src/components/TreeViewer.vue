<template>
  <TreeTable :value="nodes"
  class="p-treetable-sm"
  style="margin-bottom: 2rem"
  responsiveLayout="scroll"
  :resizableColumns="true"
  columnResizeMode="fit"
  showGridlines
  :filters="tableFilters"
  @filter="expandFiltered"
  filterMode="strict"
  v-model:expandedKeys="expandedKeys"
  >
    <Column field="name" header="Name" :expander="true"
      filterMatchMode="contains"
      >
      <template #filter>
        <InputText type="text" v-model="tableFilters['name']" class="p-column-filter" placeholder="Filter by name" />
      </template>
      <template #body="slotProps">
        <span>{{slotProps.node.data.name}}</span>

      </template>
    </Column>

    <Column field="format" header="Preview" headerStyle="width: 5%"  filterMatchMode="contains">
      <template #filter>
        <InputText type="text" v-model="tableFilters['format']" class="p-column-filter" placeholder="Filter by type" />
      </template>
      <template #body="slotProps">
        <Image v-if="slotProps.node.data.type === 'image'"
          :src="slotProps.node.data.url" preview
          class="preview-img"/>
        <p class="format-label"><i :class="slotProps.node.data.icon"></i>
        {{slotProps.node.data.format}}</p>
      </template>
    </Column>

    <Column header="Actions" headerStyle="width: 10%">
      <template #filter>
        <span class="p-buttonset">
        <Button @click="expandAll" title="Expand All" icon="fa fa-expand-arrows-alt" class="p-button-sm" />
        <Button @click="collapseAll" title="Collapse All" icon="fa fa-compress-arrows-alt" class="p-button-sm" />

        </span>
      </template>
      <template #body="slotProps">
        <a v-if="slotProps.node.data.url" target="_blank" :href="slotProps.node.data.url">
          <Button icon="p-button-small pi pi-download"/>
        </a>
      </template>
    </Column>
<!----
    <template #default="slotProps">
        <span>{{slotProps.node.label}}</span>
    </template>
    <template #url="slotProps">
        <a :href="slotProps.node.data">{{slotProps.node.label}}</a>
    </template>
    <template #image="slotProps">
      <div class="image-node">
        <span>{{slotProps.node.label}}</span>
        <img :src="slotProps.node.data"/>

      </div>
    </template>
  -->
  </TreeTable>
</template>
<script setup>
import { ref } from "vue"

 const nodes = ref(props.treeData)
 const tableFilters = ref({})
 const expandedKeys = ref({})

  const props = defineProps({
    treeData: {
      type: Object,
      default() {return {}}
   }
  })

function expandFiltered(event) {
  if (event.originalEvent.type != "input"){
    // This event is also triggered when expanding a node while a filter is applied.
    // Dont mess with that
    return
  }
  if( Object.values(event.filters).join().length <= 2 ) {
    return
  }
  expandedKeys.value = {}
  for (let node of event.filteredValue) {
    expandNode(node);
  }
}

function expandAll(){
  for (let node of nodes.value) {
    expandNode(node);
  }
}
function collapseAll(){
  expandedKeys.value = {}
}

function expandNode(node) {
  if (node.children && node.children.length) {
    expandedKeys.value[node.key] = true;

    for (let child of node.children) {
      expandNode(child);
    }
  }
}

</script>
<style lang="scss">
.p-treetable-tbody {
  .format-label {
    i.fa {
      padding-right: 5px;
    }
  }
  .preview-img>img {
    max-width: 100%;
  }
}
.p-filter-column .p-buttonset button {
  margin-top: 3px;
}
</style>