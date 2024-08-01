<template>
  <div id="app">
    <h1>--Hospital operation--</h1>

    <!-- createHospital chart -->
    <form @submit.prevent="hospital_submitHospitalCreateForm">
      <h2>Create Hospital</h2>
      <input v-model="hospital_createHospitalData.name" placeholder="Enter hospital name" />
      <input v-model="hospital_createHospitalData.contact" placeholder="Enter hospital contact" />
      <button type="submit">Submit</button>
    </form>
    <p>{{ hospital_createHospital_response }}</p>

    <!-- getHospital chart -->
    <form @submit.prevent="hospital_submitHospitalGetForm">
      <h2>Get Hospital</h2>
      <!--<input v-model="getData.id" placeholder="Enter hospital ID" />-->
      <button type="submit">Submit</button>
    </form>
    <p>{{ hospital_getHospital_response }}</p>

    <!-- getPatients chart -->
    <form @submit.prevent="hospital_submitPatientGetForm">
      <h2>Get Patients</h2>
      <input v-model="hospital_getPatientData.hospitalName" placeholder="Enter hospital name" />
      <button type="submit">Submit</button>
    </form>
    <p>{{ hospital_getPatient_response }}</p>

    <!-- modifyReport chart -->
    <form @submit.prevent="hospital_submitModifyReportForm">
      <h2>Modify Report</h2>
      <input v-model="hospital_modifyReportData.hospitalName" placeholder="Enter hospital name" />
      <input v-model="hospital_modifyReportData.patientName" placeholder="Enter patient name" />
      <input v-model="hospital_modifyReportData.symptoms" placeholder="Enter symptoms" />
      <input v-model="hospital_modifyReportData.neededDrugs" placeholder="Enter needed drugs" />
      <button type="submit">Submit</button>
    </form>
    <p>{{ hospital_modifyReport_response }}</p>

    <!-- buyDrug chart -->
    <form @submit.prevent="hospital_submitBuyDrugForm">
      <h2>Buy Drug</h2>
      <input v-model="hospital_buyDrugData.hospitalName" placeholder="Enter hospital name" />
      <input v-model="hospital_buyDrugData.manufacturerName" placeholder="Enter manufacturer name" />
      <input v-model="hospital_buyDrugData.drugName" placeholder="Enter needed drugs" />
      <button type="submit">Submit</button>
    </form>
    <p>{{ hospital_buyDrug_response }}</p>

    <!-- traceDrug chart -->
    <form @submit.prevent="hospital_submitTraceDrugForm">
      <h2>Trace Drug</h2>
      <input v-model="hospital_traceDrugData.traceCode" placeholder="Enter trace code" />
      <button type="submit">Submit</button>
    </form>
    <p>{{ hospital_traceDrug_response }}</p>


    <h1>--manufacturer operation--</h1>

    <!-- createManufacturer chart -->
    <form @submit.prevent="hospital_submitHospitalCreateForm">
      <h2>Create Manufacturer</h2>
      <input v-model="hospital_createHospitalData.name" placeholder="Enter manufacturer name" />
      <input v-model="hospital_createHospitalData.contact" placeholder="Enter manufacturer contact" />
      <button type="submit">Submit</button>
    </form>
    <p>{{ hospital_createHospital_response }}</p>

    <!-- getManufacturers chart -->
    <form @submit.prevent="hospital_submitHospitalGetForm">
      <h2>Get Manufacturer</h2>
      <!--<input v-model="getData.id" placeholder="Enter hospital ID" />-->
      <button type="submit">Submit</button>
    </form>
    <p>{{ hospital_getHospital_response }}</p>

    <!-- getManufacturerInfo chart -->
    <form @submit.prevent="hospital_submitHospitalGetForm">
      <h2>Get ManufacturerInfo</h2>
      <input v-model="hospital_getPatientData.id" placeholder="Enter Manufacturer Name" />
      <button type="submit">Submit</button>
    </form>
    <p>{{ hospital_getHospital_response }}</p>

    <!-- produceDrug chart -->
    <form @submit.prevent="hospital_submitHospitalGetForm">
      <h2>Produce Drug</h2>
      <input v-model="hospital_getPatientData.id" placeholder="Enter Manufacturer Name" />
      <input v-model="hospital_getPatientData.id" placeholder="Enter Drug Name" />
      <input v-model="hospital_getPatientData.id" placeholder="Enter Price" />
      <button type="submit">Submit</button>
    </form>
    <p>{{ hospital_getHospital_response }}</p>

    <h1>--patient operation--</h1>

    <!-- createPatient chart -->
    <form @submit.prevent="hospital_submitHospitalGetForm">
      <h2>Create Patient</h2>
      <input v-model="hospital_getPatientData.id" placeholder="Enter Name" />
      <input v-model="hospital_getPatientData.id" placeholder="Enter BirthDate" />
      <input v-model="hospital_getPatientData.id" placeholder="Enter Height" />
      <input v-model="hospital_getPatientData.id" placeholder="Enter Weight" />
      <input v-model="hospital_getPatientData.id" placeholder="Enter Gender" />
      <input v-model="hospital_getPatientData.id" placeholder="Enter Contact" />
      <button type="submit">Submit</button>
    </form>

    <!-- getPatient chart -->
    <form @submit.prevent="hospital_submitHospitalGetForm">
      <h2>Get Patients</h2>
      <button type="submit">Submit</button>
    </form>

    <!-- conductExamination chart -->
    <form @submit.prevent="hospital_submitHospitalGetForm">
      <h2>Conduct Examination</h2>
      <input v-model="hospital_getPatientData.id" placeholder="Enter Hospital Name" />
      <input v-model="hospital_getPatientData.id" placeholder="Enter Patient Name" />
      <button type="submit">Submit</button>
    </form>

    <!-- viewReport chart -->
    <form @submit.prevent="hospital_submitHospitalGetForm">
      <h2>View Report</h2>
      <input v-model="hospital_getPatientData.id" placeholder="Enter Patient Name" />
      <input v-model="hospital_getPatientData.id" placeholder="Enter Hospital Name" />
      <input v-model="hospital_getPatientData.id" placeholder="Enter ReportID" />
      <button type="submit">Submit</button>
    </form>

    <!-- buyDrug chart -->
    <form @submit.prevent="hospital_submitHospitalGetForm">
      <h2>Buy Drug</h2>
      <input v-model="hospital_getPatientData.id" placeholder="Enter Hospital Name" />
      <input v-model="hospital_getPatientData.id" placeholder="Enter Drug Name" />
      <button type="submit">Submit</button>
    </form>

    <!-- traceDrug chart -->
    <form @submit.prevent="hospital_submitHospitalGetForm">
      <h2>Trace Drug</h2>
      <input v-model="hospital_getPatientData.id" placeholder="Enter TraceID" />
      <button type="submit">Submit</button>
    </form>

  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'App',
  data() {
    return {
      hospital_createHospitalData: {
        name: '',
        contact: ''
      }, 
      hospital_getHospitalData: {
        id: ''
      },
      hospital_getPatientData: {
        hospitalName: ''
      },
      hospital_modifyReportData: {
        hospitalName: '',
        patientName: '',
        symptoms: '',
        neededDrugs: ''
      }, 
      hospital_buyDrugData: {
        hospitalName: '',
        manufacturerName: '',
        drugName: ''
      },
      hospital_traceDrugData: {
        traceCode: ''
      },
      hospital_createHospital_response: '', 
      hospital_getHospital_response: '',
      hospital_getPatient_response: '', 
      hospital_modifyReport_response: '', 
      hospital_buyDrug_response: '', 
      hospital_traceDrug_response: '', 

      manufacturer_createManufacturerlData: {
        name: '',
        contact: ''
      }, 
    };
  },
  methods: {
    async hospital_submitHospitalCreateForm() {//event when submit createHospital form
      try {
        const response = await axios.post('http://47.237.161.142:8888/contract/hospital/createHospital', this.hospital_createHospitalData);
        if (response.data.code === 1) {
          this.hospital_createHospital_response = response.data.message; 
        } else {
          this.hospital_createHospital_response = "Submitted successfully";
        }
      } catch (error) {
        console.error("There was an error!", error);
        this.hospital_createHospital_response = "Error submitting form.";
      }
    },
    async hospital_submitHospitalGetForm() {//event when submit getHospital form
      try {
        const response = await axios.post('http://47.237.161.142:8888/contract/hospital/getHospitals', this.hospital_getHospitalData);
        if (response.data.code === 1) {
          this.hospital_getHospital_response = response.data.message; 
          this.hospital_getHospital_response = "Retrieved successfully: " + response.data.data;
        }
      } catch (error) {
        console.error("There was an error!", error);
        this.hospital_getHospital_response = "Error submitting form.";
      }
    },
    async hospital_submitPatientGetForm() {//event when submit getPatients form
      try {
        const response = await axios.post('http://47.237.161.142:8888/contract/hospital/getPatients', this.hospital_getPatientData);
        if (response.data.code === 1) {
          this.hospital_getPatient_response = response.data.message;
        } else {
          this.hospital_getPatient_response = "Retrieved successfully: " + response.data.data;
        }
      } catch (error) {
        console.error("There was an error!", error);
        this.hospital_getPatient_response = "Error submitting form.";
      }
    },
    async hospital_submitModifyReportForm() {//event when submit modifyReport form
      try {
        const response = await axios.post('http://47.237.161.142:8888/contract/hospital/modifyReport', this.hospital_modifyReportData);
        if (response.data.code === 1) {
          this.hospital_modifyReport_response = response.data.message; 
        } else {
          this.hospital_modifyReport_response = "Retrieved successfully: " + response.data.data;
        }
      } catch (error) {
        console.error("There was an error!", error);
        this.hospital_modifyReport_response = "Error submitting form.";
      }
    },
    async hospital_submitBuyDrugForm() {//event when submit buyDrug form
      try {
        const response = await axios.post('http://47.237.161.142:8888/contract/hospital/buyDrug', this.hospital_buyDrugData);
        if (response.data.code === 1) {
          this.hospital_buyDrug_response = response.data.message; 
        } else {
          this.hospital_buyDrug_response = "Retrieved successfully: " + response.data.data;
        }
      } catch (error) {
        console.error("There was an error!", error);
        this.hospital_buyDrug_response = "Error submitting form.";
      }
    },
    async hospital_submitTraceDrugForm() {//event when submit traceDrug form
      try {
        const response = await axios.post('http://47.237.161.142:8888/contract/hospital/traceDrug', this.hospital_traceDrugData);
        if (response.data.code === 1) {
          this.hospital_traceDrug_response = "Invalid TraceCode";
        } else {
          this.hospital_traceDrug_response = "Retrieved successfully: " + response.data.message;
        }
      } catch (error) {
        console.error("There was an error!", error);
        this.hospital_traceDrug_response = "Error submitting form.";
      }
    }
  }
};
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

h1 {
  font-size: 2em;
}

h2 {
  font-size: 1.5em;
  margin-top: 20px;
}

form {
  margin: 20px 0;
}

input {
  margin-right: 10px;
}

button {
  padding: 5px 10px;
}

p {
  font-size: 1.2em;
  color: green;
}
</style>
