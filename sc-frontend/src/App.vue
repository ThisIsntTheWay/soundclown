<template>
  <div id="app" class="container">
    <div class="row">
      <div class="col-md-6 offset-md-3 py-5">
        <h1>Download SoundCloud track</h1>
        Base URL: {{ baseUrl }}

        <form v-on:submit.prevent="downloadTracks">
          <div class="form-group">
            <input v-model="targetUrl" type="text" placeholder="Enter a soundcloud URL" class="form-control">
          </div>
          <div class="form-group">
            <button class="btn btn-primary">Let's get it</button>
            <p></p>
            Response: {{ apiResponse }}
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import { response } from 'express';

export default {
  name: 'App',

  data() {
    return {
      targetUrl: '',
      baseUrl: window.location.protocol + "//" + window.location.host,
      apiResponse: null,
    }
  },

  methods: {
    downloadTracks(targetUrl) {
      axios.post(this.baseUrl + "/api/download", {
        url:  targetUrl
      })
      .then((response) => {
        this.apiResponse = response.data;
      })
      .catch((error) => {
        this.apiResponse = "ERR: ", error, " - ", response.data.error;
        //window.alert(`The API returned an error: ${error}`);
      })
    }
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
