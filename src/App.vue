<template>
  <div id="app">
    <h1>Check-It</h1>
    <checklist/>
  </div>
</template>

<script>
import Checklist from './components/Checklist'

export default {
  name: 'App',
  components: {
    Checklist
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  /* color: #2c3e50; */
  color: green;
}

span {
  color: green;
}

button {
  background-color: white;
  color: green;
  border-radius: 50%;
  font-size: 16px;
}

input:focus,
select:focus,
textarea:focus,
button:focus {
  outline: none;
}
</style>
