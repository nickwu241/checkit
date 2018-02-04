<template>
  <div>
    <p>List Name:
      <input type="text" v-model="listname">
    </p>
    <transition name="fade">
    <h5 v-if="sharing">{{ API_URL + '?name=' + listname }}</h5>
    </transition>
    <vue-button-spinner
      :isLoading="saving"
      :disabled="saving"
      :status="savingStatus">
    <span @click="saveList()">Save</span>
    </vue-button-spinner>
    <vue-button-spinner
      :isLoading="false"
      :disabled="sharing"
      :status="sharingStatus">
    <span @click="share()">Share</span>
    </vue-button-spinner>
    <br>
    <br>
    <br>
    <div v-for="(item, i) in checklist" v-bind:key="item.id">
      <input id="checkBox" type="checkbox" :checked="item.checked" @click="clicked(i)">
      <input type="text" v-model="item.text"
        @input="onChange()"
        @change="onChange()"
        @paste="onChange()"
        @keyup.enter="addItem()">
      <button @click="deleteItem(i)">X</button>
    </div>
  <button @click="addItem()">Add</button>
  <br>
  </div>
</template>

<script>
import VueButtonSpinner from 'vue-button-spinner'

var id = null
if (document.URL.includes('/checklist')) {
  id = document.URL.substring(document.URL.indexOf('=') + 1)
  console.log('ID:', id)
}
// const = 'http://localhost:8080'
const BASE_URL = 'http://3d5a79cd.ngrok.io'
export default {
  name: 'Checklist',
  components: {
    VueButtonSpinner
  },
  data () {
    return {
      firstId: id,
      listname: id || 'default',
      API_URL: BASE_URL + '/checklist',
      checklist: [],
      shouldUpdate: true,
      sharing: false,
      sharingStatus: '',
      saving: false,
      savingStatus: ''
    }
  },
  methods: {
    firstFetchItems: function () {
      if (this.firstId) {
        fetch(this.API_URL + '?name=' + id).then(r => {
          return r.json()
        }).then(r => {
          if (this.shouldUpdate) {
            this.checklist = r.items
          }
        })
      } else {
        this.fetchItems()
        // fetch(this.API_URL + '?default=true').then(r => {
        //   return r.json()
        // }).then(r => {
        //   if (this.shouldUpdate) {
        //     this.checklist = r.items
        //   }
        // })
      }
    },
    fetchItems: function () {
      fetch(this.API_URL).then(r => {
        return r.json()
      }).then(r => {
        if (this.shouldUpdate) {
          this.checklist = r.items
        }
      })
    },
    postItems: function () {
      this.shouldUpdate = false

      const data = JSON.stringify({'checklist': this.checklist})
      console.log(data)
      fetch(this.API_URL, {
        method: 'POST',
        body: data
      }).then(r => {
        console.log('OK', r.ok)
        this.shouldUpdate = true
      })
    },
    saveList: function () {
      this.saving = true
      const data = JSON.stringify({'name': this.listname})

      fetch(this.API_URL + '/share', {
        method: 'POST',
        body: data
      }).then(r => {
        console.log('OK', r.ok)
        this.saving = false
        this.savingStatus = true
        setTimeout(() => { this.savingStatus = '' }, 2000)
      }).catch(err => {
        console.error(err)
        this.isLoading = false
        this.status = false
      })
    },
    clicked: function (i) {
      this.checklist[i].checked = !this.checklist[i].checked
      this.postItems()
    },
    addItem: function () {
      this.checklist.push({checked: false, text: ''})
      this.postItems()
    },
    deleteItem: function (i) {
      if (this.checklist.length === 1) {
        this.checklist = [{'checked': false, 'text': ''}]
        this.postItems()
        return
      }
      this.checklist.splice(i, 1)
      this.postItems()
    },
    share: function () {
      this.sharing = true
      setTimeout(() => {
        this.sharing = false
        this.sharingStatus = ''
      }, 7000)
    },
    onChange: function () {
      this.postItems()
    }
  },
  mounted: function () {
    this.firstFetchItems()
    setInterval(function () {
      if (this.shouldUpdate) {
        this.fetchItems()
      }
    }.bind(this), 1000)
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.fade-enter-active, .fade-leave-active {
  transition: opacity .5s;
}
.fade-enter, .fade-leave-to {
  opacity: 0;
}
</style>
