<template>
  <div class="hello">
    <!-- <h1>{{ msg }}</h1> -->

    <div>
      <h1 class="title">Sizes</h1>
    </div>
    <hr>

    <div>
      <b-container fluid>
        <b-row>
          <b-col v-for="size in sizes" :key="size">
            <b-card  align="center" :title="'Pack Size:' + size" class="mb-2">

              <b-button v-on:click="deletePack(size)" variant="danger">Delete</b-button>
            </b-card>
          </b-col>
        </b-row>
        <b-row cols-sm="6" align-h="center">
          <b-col align-self="center">
            <b-input-group size="md" prepend="Pack Size">
              <b-form-input v-model.number="newSize" type="number"></b-form-input>
            </b-input-group>
            <b-button variant="primary" v-on:click="writePack()">Create</b-button>
          </b-col>
        </b-row>
      </b-container>
    </div>

    <hr>

    <div>
      <h1 class="title">Packs</h1>
    </div>
    <hr>

    <div>
      <b-container fluid>
        <b-row cols-sm="6" align-h="center">
          <b-col align-self="center">
            <b-input-group size="md" prepend="Target">
              <b-form-input v-model.number="target" type="number"></b-form-input>
            </b-input-group>
            <b-button variant="primary" v-on:click="orderBest()">Best</b-button>
            <b-button variant="primary" v-on:click="orderFast()">Fast</b-button>
          </b-col>
        </b-row>
        <b-row>
          <b-col v-for="pack in packs" :key="pack.size">
            <b-card align="center" header="Pack" class="mb-2">
              <b-card-text>Count: {{pack.count}}</b-card-text>
              <b-card-text>Size: {{pack.size}}</b-card-text>
            </b-card>
          </b-col>
        </b-row>
      </b-container>
    </div>

    <div>
      <b-modal ref="error-modal" id="error-modal" centered title="Error">
        <p class="my-4">{{this.error}}</p>
      </b-modal>
    </div>

  </div>
</template>

<script>

import axios from 'axios';

export default {
  name: 'Calculator',

  data: function() {
    return {
      sizes: [],
      newSize: 0,
      target: 0,
      packs: [],
      error: "",
    }
  },
  created() {
    this.readPack()
  },
  methods: {
    showModal() {
      this.$refs['error-modal'].show()
    },
    readPack: function() {
      axios({ method: "GET", url: "http://"+ process.env.VUE_APP_ENDPOINT + "/packs/read", data: null, headers: {"content-type": "text/plain" } }).then(result => {
        this.sizes = result.data

      }).catch( error => {
        this.error = error.response.data.error
        this.showModal()
      });
    },
    deletePack: function(size) {
      axios({ method: "DELETE", url: "http://"+ process.env.VUE_APP_ENDPOINT + "/packs/delete", data: {"size": size}, headers: {"content-type": "text/plain" } }).then(result => {
        console.log(result)
        this.readPack()

      }).catch( error => {
        this.error = error.response.data.error
        this.showModal()
      });
    },
    writePack: function() {
      axios({ method: "POST", url: "http://"+ process.env.VUE_APP_ENDPOINT + "/packs/write", data: {"size":this.newSize}, headers: {"content-type": "text/plain" } }).then(result => {
        console.log(result)
        this.readPack()

      }).catch( error => {
        this.error = error.response.data.error
        this.showModal()
      });
    },
    orderBest: function() {
      axios({ method: "POST", url: "http://"+ process.env.VUE_APP_ENDPOINT + "/order/best", data: {"target":this.target}, headers: {"content-type": "text/plain" } }).then(result => {
        this.packs = result.data

      }).catch( error => {
        this.error = error.response.data.error
        this.showModal()
      });
    },
    orderFast: function() {
      axios({ method: "POST", url: "http://"+ process.env.VUE_APP_ENDPOINT + "/order/fast", data: {"target":this.target}, headers: {"content-type": "text/plain" } }).then(result => {
        this.packs = result.data

      }).catch( error => {
        this.error = error.response.data.error
        this.showModal()
      });
    },
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>