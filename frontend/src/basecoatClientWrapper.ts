import Cookies from "js-cookie";
import { BasecoatClient } from "./BasecoatServiceClientPb";
import {
  Address,
  Base,
  Colorant,
  Contact,
  Contractor,
  Formula,
  Job,
} from "./basecoat_message_pb";
import {
  CreateAPITokenRequest,
  CreateContractorRequest,
  CreateFormulaRequest,
  CreateJobRequest,
  DeleteContractorRequest,
  DeleteFormulaRequest,
  DeleteJobRequest,
  GetContractorRequest,
  GetFormulaRequest,
  GetJobRequest,
  GetSystemInfoRequest,
  ListContractorsRequest,
  ListFormulasRequest,
  ListJobsRequest,
  SearchFormulasRequest,
  SearchJobsRequest,
  UpdateContractorRequest,
  UpdateFormulaRequest,
  UpdateJobRequest,
} from "./basecoat_transport_pb";

interface LoginInfo {
  username: string;
  password: string;
}

interface formulaMap {
  [key: string]: Formula;
}
interface jobMap {
  [key: string]: Job;
}

interface contractorMap {
  [key: string]: Contractor;
}

interface systemInfo {
  build_time: string;
  commit: string;
  debug_enabled: boolean;
  frontend_enabled: boolean;
  semver: string;
}

// BasecoatClientWrapper is a wrapper for all frontend to backend communication
class BasecoatClientWrapper {
  client: BasecoatClient;

  constructor() {
    let url = location.protocol + "//" + location.host;
    this.client = new BasecoatClient(url, null, null);
  }

  // isUserLoggedIn determines if the user should be kicked back to the login route
  isUserLoggedIn(): boolean {
    if (!Cookies.get("username") || !Cookies.get("token")) {
      return false;
    }
    return true;
  }

  // handleLogin is a Promise that returns whether a login was successful
  // or not and sets relevant cookies
  handleLogin(loginInfo: LoginInfo): Promise<string> {
    let tokenRequest = new CreateAPITokenRequest();
    tokenRequest.setUser(loginInfo.username);
    tokenRequest.setPassword(loginInfo.password);
    tokenRequest.setDuration(10368000); // Four months

    return new Promise((resolve, reject) => {
      this.client.createAPIToken(tokenRequest, {}, function(err, response) {
        if (err) {
          reject(err);
          return;
        }
        Cookies.set("username", loginInfo.username, {
          expires: 120,
          secure: true,
        });
        Cookies.set("token", response.getKey(), { expires: 120, secure: true });
        resolve();
      });
    });
  }

  // handleLogout signs a user out by removing relevant cookies
  handleLogout() {
    Cookies.remove("username");
    Cookies.remove("token");
  }

  //getFormula retrieves a single formula by ID
  getFormula(formulaID: string): Promise<Formula | undefined> {
    let metadata = { Authorization: "Bearer " + Cookies.get("token") };
    let getFormulaRequest = new GetFormulaRequest();
    getFormulaRequest.setId(formulaID);

    return new Promise((resolve, reject) => {
      this.client.getFormula(getFormulaRequest, metadata, function(
        err,
        response
      ) {
        if (err) {
          reject(err);
          return;
        }
        resolve(response.getFormula());
      });
    });
  }

  //getFormulaData retrieves all formulas from the backend
  getFormulaData(): Promise<formulaMap | undefined> {
    let metadata = { Authorization: "Bearer " + Cookies.get("token") };
    let listFormulasRequest = new ListFormulasRequest();

    return new Promise((resolve, reject) => {
      if (!this.isUserLoggedIn()) {
        reject(undefined);
      }
      this.client.listFormulas(listFormulasRequest, metadata, function(
        err,
        response
      ) {
        if (err) {
          reject(err);
        }

        let formulas: formulaMap = {};
        response.getFormulasMap().forEach(function(value, key) {
          formulas[key] = value;
        });

        resolve(formulas);
      });
    });
  }

  //searchFormulas returns formulas that match the search term
  searchFormulas(searchTerm: string): Promise<string[] | undefined> {
    let metadata = { Authorization: "Bearer " + Cookies.get("token") };
    let searchFormulasRequest = new SearchFormulasRequest();
    searchFormulasRequest.setTerm(searchTerm);

    return new Promise((resolve, reject) => {
      this.client.searchFormulas(searchFormulasRequest, metadata, function(
        err,
        response
      ) {
        if (err) {
          reject(err);
          return;
        }
        resolve(response.getResultsList());
      });
    });
  }

  //getJob retrieves a single job by ID
  getJob(jobID: string): Promise<Job | undefined> {
    let metadata = { Authorization: "Bearer " + Cookies.get("token") };
    let getJobRequest = new GetJobRequest();
    getJobRequest.setId(jobID);

    return new Promise((resolve, reject) => {
      this.client.getJob(getJobRequest, metadata, function(err, response) {
        if (err) {
          reject(err);
          return;
        }
        resolve(response.getJob());
      });
    });
  }

  //getJobData retrieves all jobs from the backend
  getJobData(): Promise<jobMap | undefined> {
    let metadata = { Authorization: "Bearer " + Cookies.get("token") };
    let listJobsRequest = new ListJobsRequest();

    return new Promise((resolve, reject) => {
      if (!this.isUserLoggedIn()) {
        reject(undefined);
      }
      this.client.listJobs(listJobsRequest, metadata, function(err, response) {
        if (err) {
          reject(err);
        }

        let jobs: jobMap = {};
        response.getJobsMap().forEach(function(value, key) {
          jobs[key] = value;
        });

        resolve(jobs);
      });
    });
  }

  //searchJobs returns jobs that match the search term
  searchJobs(searchTerm: string): Promise<string[] | undefined> {
    let metadata = { Authorization: "Bearer " + Cookies.get("token") };
    let searchJobsRequest = new SearchJobsRequest();
    searchJobsRequest.setTerm(searchTerm);

    return new Promise((resolve, reject) => {
      this.client.searchJobs(searchJobsRequest, metadata, function(
        err,
        response
      ) {
        if (err) {
          reject(err);
          return;
        }
        resolve(response.getResultsList());
      });
    });
  }

  //submitCreateFormulaForm submits the formula create form
  submitCreateFormulaForm(
    formulaData: CreateFormulaRequest.AsObject
  ): Promise<string> {
    let metadata = { Authorization: "Bearer " + Cookies.get("token") };
    return new Promise((resolve, reject) => {
      if (!this.isUserLoggedIn()) {
        reject();
      }

      let createFormulaRequest = new CreateFormulaRequest();
      createFormulaRequest.setName(formulaData.name);
      createFormulaRequest.setNumber(formulaData.number);
      createFormulaRequest.setNotes(formulaData.notes);
      createFormulaRequest.setJobsList(formulaData.jobsList);

      let basesList: Base[] = [];
      formulaData.basesList.forEach(function(item, index) {
        let newBase = new Base();
        newBase.setType(item.type);
        newBase.setName(item.name);
        newBase.setAmount(item.amount);

        basesList.push(newBase);
      });

      createFormulaRequest.setBasesList(basesList);

      let colorantsList: Colorant[] = [];
      formulaData.colorantsList.forEach(function(item, index) {
        let newColorant = new Colorant();
        newColorant.setType(item.type);
        newColorant.setName(item.name);
        newColorant.setAmount(item.amount);

        colorantsList.push(newColorant);
      });
      createFormulaRequest.setColorantsList(colorantsList);

      this.client.createFormula(createFormulaRequest, metadata, function(
        err,
        response
      ) {
        if (err) {
          reject(err);
        }
        resolve();
      });
    });
  }

  //submitCreateJobForm submits a new job to the backend
  submitCreateJobForm(jobData: CreateJobRequest.AsObject): Promise<string> {
    let metadata = { Authorization: "Bearer " + Cookies.get("token") };
    return new Promise((resolve, reject) => {
      if (!this.isUserLoggedIn()) {
        reject();
      }

      let createJobRequest = new CreateJobRequest();
      createJobRequest.setName(jobData.name);

      if (jobData.address != undefined) {
        let address = new Address();
        address.setStreet(jobData.address.street);
        address.setStreet2(jobData.address.street2);
        address.setCity(jobData.address.city);
        address.setState(jobData.address.state);
        address.setZipcode(jobData.address.zipcode);
        createJobRequest.setAddress(address);
      }

      createJobRequest.setNotes(jobData.notes);
      createJobRequest.setFormulasList(jobData.formulasList);

      if (jobData.contact != undefined) {
        let contact = new Contact();
        contact.setName(jobData.contact.name);
        contact.setEmail(jobData.contact.email);
        contact.setPhone(jobData.contact.phone);
        createJobRequest.setContact(contact);
      }

      createJobRequest.setContractorId(jobData.contractorId);

      this.client.createJob(createJobRequest, metadata, function(
        err,
        response
      ) {
        if (err) {
          reject(err);
        }
        resolve(response.getJob()?.getId());
      });
    });
  }

  submitManageFormulaForm(
    formulaData: UpdateFormulaRequest.AsObject
  ): Promise<string> {
    let metadata = { Authorization: "Bearer " + Cookies.get("token") };
    return new Promise((resolve, reject) => {
      if (!this.isUserLoggedIn()) {
        reject();
      }

      let updateFormulaRequest = new UpdateFormulaRequest();
      updateFormulaRequest.setId(formulaData.id);
      updateFormulaRequest.setName(formulaData.name);
      updateFormulaRequest.setNumber(formulaData.number);
      updateFormulaRequest.setNotes(formulaData.notes);
      updateFormulaRequest.setJobsList(formulaData.jobsList);

      let basesList: Base[] = [];
      formulaData.basesList.forEach(function(item, index) {
        let newBase = new Base();
        newBase.setType(item.type);
        newBase.setName(item.name);
        newBase.setAmount(item.amount);

        basesList.push(newBase);
      });

      updateFormulaRequest.setBasesList(basesList);

      let colorantsList: Colorant[] = [];
      formulaData.colorantsList.forEach(function(item, index) {
        let newColorant = new Colorant();
        newColorant.setType(item.type);
        newColorant.setName(item.name);
        newColorant.setAmount(item.amount);

        colorantsList.push(newColorant);
      });
      updateFormulaRequest.setColorantsList(colorantsList);

      this.client.updateFormula(updateFormulaRequest, metadata, function(
        err,
        response
      ) {
        if (err) {
          reject(err);
        }
        resolve();
      });
    });
  }

  submitManageJobForm(jobData: UpdateJobRequest.AsObject): Promise<string> {
    let metadata = { Authorization: "Bearer " + Cookies.get("token") };
    return new Promise((resolve, reject) => {
      if (!this.isUserLoggedIn()) {
        reject();
      }

      let updateJobRequest = new UpdateJobRequest();
      updateJobRequest.setId(jobData.id);
      updateJobRequest.setName(jobData.name);
      if (jobData.address != undefined) {
        let address = new Address();
        address.setStreet(jobData.address.street);
        address.setStreet2(jobData.address.street2);
        address.setCity(jobData.address.city);
        address.setState(jobData.address.state);
        address.setZipcode(jobData.address.zipcode);
        updateJobRequest.setAddress(address);
      }

      updateJobRequest.setNotes(jobData.notes);
      updateJobRequest.setFormulasList(jobData.formulasList);

      if (jobData.contact != undefined) {
        let contact = new Contact();
        contact.setName(jobData.contact.name);
        contact.setEmail(jobData.contact.email);
        contact.setPhone(jobData.contact.phone);
        updateJobRequest.setContact(contact);
      }

      updateJobRequest.setContractorId(jobData.contractorId);

      this.client.updateJob(updateJobRequest, metadata, function(
        err,
        response
      ) {
        if (err) {
          reject(err);
        }
        resolve();
      });
    });
  }

  deleteFormula(formulaID: string): Promise<string> {
    let metadata = { Authorization: "Bearer " + Cookies.get("token") };
    let deleteFormulaRequest = new DeleteFormulaRequest();
    deleteFormulaRequest.setId(formulaID);

    return new Promise((resolve, reject) => {
      if (!this.isUserLoggedIn()) {
        reject();
      }

      this.client.deleteFormula(deleteFormulaRequest, metadata, function(
        err,
        response
      ) {
        if (err) {
          reject(err);
        }
        resolve();
      });
    });
  }

  deleteJob(jobID: string): Promise<string> {
    let metadata = { Authorization: "Bearer " + Cookies.get("token") };
    let deleteJobRequest = new DeleteJobRequest();
    deleteJobRequest.setId(jobID);

    return new Promise((resolve, reject) => {
      if (!this.isUserLoggedIn()) {
        reject();
      }

      this.client.deleteJob(deleteJobRequest, metadata, function(
        err,
        response
      ) {
        if (err) {
          reject(err);
        }
        resolve();
      });
    });
  }

  //getContractor retrieves a single contractor by ID
  getContractor(contractorID: string): Promise<Contractor | undefined> {
    let metadata = { Authorization: "Bearer " + Cookies.get("token") };
    let getContractorRequest = new GetContractorRequest();
    getContractorRequest.setId(contractorID);

    return new Promise((resolve, reject) => {
      this.client.getContractor(getContractorRequest, metadata, function(
        err,
        response
      ) {
        if (err) {
          reject(err);
          return;
        }
        resolve(response.getContractor());
      });
    });
  }

  //getContractorData retrieves all contractors from the backend
  getContractorData(): Promise<contractorMap | undefined> {
    let metadata = { Authorization: "Bearer " + Cookies.get("token") };
    let listContractorsRequest = new ListContractorsRequest();

    return new Promise((resolve, reject) => {
      if (!this.isUserLoggedIn()) {
        reject(undefined);
      }
      this.client.listContractors(listContractorsRequest, metadata, function(
        err,
        response
      ) {
        if (err) {
          reject(err);
        }

        let contractors: contractorMap = {};
        response.getContractorsMap().forEach(function(value, key) {
          contractors[key] = value;
        });

        resolve(contractors);
      });
    });
  }

  //submitCreateContractorForm submits a new contractor to the backend
  submitCreateContractorForm(
    contractorData: CreateContractorRequest.AsObject
  ): Promise<string> {
    let metadata = { Authorization: "Bearer " + Cookies.get("token") };
    return new Promise((resolve, reject) => {
      if (!this.isUserLoggedIn()) {
        reject();
      }

      let createContractorRequest = new CreateContractorRequest();
      createContractorRequest.setCompany(contractorData.company);

      if (contractorData.contact != undefined) {
        let contact = new Contact();
        contact.setName(contractorData.contact.name);
        contact.setEmail(contractorData.contact.email);
        contact.setPhone(contractorData.contact.phone);
        createContractorRequest.setContact(contact);
      }

      createContractorRequest.setJobsList(contractorData.jobsList);

      this.client.createContractor(createContractorRequest, metadata, function(
        err,
        response
      ) {
        if (err) {
          reject(err);
        }
        resolve(response.getContractor()?.getId());
      });
    });
  }

  submitManageContractorForm(
    contractorData: UpdateContractorRequest.AsObject
  ): Promise<string> {
    let metadata = { Authorization: "Bearer " + Cookies.get("token") };
    return new Promise((resolve, reject) => {
      if (!this.isUserLoggedIn()) {
        reject();
      }

      let updateContractorRequest = new UpdateContractorRequest();
      updateContractorRequest.setId(contractorData.id);
      updateContractorRequest.setCompany(contractorData.company);

      if (contractorData.contact != undefined) {
        let contact = new Contact();
        contact.setName(contractorData.contact.name);
        contact.setEmail(contractorData.contact.email);
        contact.setPhone(contractorData.contact.phone);
        updateContractorRequest.setContact(contact);
      }

      updateContractorRequest.setJobsList(contractorData.jobsList);

      this.client.updateContractor(updateContractorRequest, metadata, function(
        err,
        response
      ) {
        if (err) {
          reject(err);
        }
        resolve();
      });
    });
  }

  deleteContractor(contractorID: string): Promise<string> {
    let metadata = { Authorization: "Bearer " + Cookies.get("token") };
    let deleteContractorRequest = new DeleteContractorRequest();
    deleteContractorRequest.setId(contractorID);

    return new Promise((resolve, reject) => {
      if (!this.isUserLoggedIn()) {
        reject();
      }

      this.client.deleteContractor(deleteContractorRequest, metadata, function(
        err,
        response
      ) {
        if (err) {
          reject(err);
        }
        resolve();
      });
    });
  }

  //getSystemInfo retrieves a system information
  getSystemInfo(): Promise<systemInfo | undefined> {
    let getSystemInfoRequest = new GetSystemInfoRequest();

    return new Promise((resolve, reject) => {
      this.client.getSystemInfo(getSystemInfoRequest, {}, function(
        err,
        response
      ) {
        if (err) {
          reject(err);
          return;
        }
        let systemInfo: systemInfo = {
          build_time: response.getBuildTime(),
          commit: response.getCommit(),
          debug_enabled: response.getDebugEnabled(),
          frontend_enabled: response.getFrontendEnabled(),
          semver: response.getSemver(),
        };
        resolve(systemInfo);
      });
    });
  }
}

export default BasecoatClientWrapper;
