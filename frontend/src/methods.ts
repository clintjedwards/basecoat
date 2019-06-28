import Vue from 'vue'
import VueCookies from 'vue-cookies'
import store from './store';
import * as bcInterface from './basecoatInterfaces'
import { app } from './index'

import {
    ListFormulasRequest,
    Formula, Job,
    GetFormulaRequest,
    ListJobsRequest,
    CreateAPITokenRequest,
    CreateFormulaRequest,
    Colorant,
    Base
} from "./basecoat_pb"
import { BasecoatClient } from './BasecoatServiceClientPb';

Vue.use(VueCookies)

export class BasecoatFrontend {
    client: BasecoatClient;
    constructor(client: BasecoatClient) {
        this.client = client;
    }

    VerifyLogin() {
        if (!Vue.cookies.isKey('username') || !Vue.cookies.isKey('token')) {
            store.commit('updateLoginState', false)
            return
        }

        store.commit('updateUsername', Vue.cookies.get('username'))
        store.commit('updateLoginState', true)
        this.LoadFormulaData();
        this.LoadJobData();
    }

    HandleLogin(loginInfo: bcInterface.LoginInfo) {
        let self = this
        store.commit('updateLoginIsLoading', true)

        let tokenRequest = new CreateAPITokenRequest();
        tokenRequest.setUser(loginInfo.username);
        tokenRequest.setPassword(loginInfo.password);
        tokenRequest.setDuration(10368000);
        this.client.createAPIToken(tokenRequest, {}, function (err, response) {
            if (err) {
                console.log(err)
                store.commit('displaySnackBar', "Invalid Login Credentials")
                store.commit('updateLoginState', false)
                store.commit('updateLoginIsLoading', false)
                return
            }
            Vue.cookies.set('username', loginInfo.username, "4m", undefined, undefined, true)
            Vue.cookies.set('token', response.getKey(), "4m", undefined, undefined, true)
            store.commit('updateUsername', loginInfo.username)
            store.commit('updateLoginState', true)
            store.commit('updateLoginIsLoading', false)
            self.LoadFormulaData()
            self.LoadJobData()
        })
    }

    HandleLogout() {
        Vue.cookies.remove('username')
        Vue.cookies.remove('token')
        this.VerifyLogin()
    }

    GetFormula(formulaID: string): Formula | undefined {
        let getFormulaRequest = new GetFormulaRequest();
        getFormulaRequest.setId(formulaID)
        let metadata = { 'Authorization': 'Bearer ' + Vue.cookies.get('token') }
        this.client.getFormula(getFormulaRequest, metadata, function (err, response) {
            if (err) {
                console.log(err)
                store.commit('displaySnackBar', "Could not load formula.")
                return undefined
            }

            let formula = response.getFormula();
            return formula
        })

        return undefined
    }

    LoadFormulaData() {
        let listFormulasRequest = new ListFormulasRequest();
        let metadata = { 'Authorization': 'Bearer ' + Vue.cookies.get('token') }
        this.client.listFormulas(listFormulasRequest, metadata, function (err, response) {
            if (err) {
                console.log(err)
                store.commit('displaySnackBar', "Could not load current formulas.")
                return
            }

            let formulaMap: { [key: string]: Formula } = {}
            response.getFormulasMap().forEach(function (value, key) {
                formulaMap[key] = value
            })
            store.commit('updateFormulaData', formulaMap)
            store.commit('updateTotalFormulas')
        })
    }

    LoadJobData() {
        let listJobsRequest = new ListJobsRequest();
        let metadata = { 'Authorization': 'Bearer ' + Vue.cookies.get('token') }
        this.client.listJobs(listJobsRequest, metadata, function (err, response) {
            if (err) {
                console.log(err)
                store.commit('displaySnackBar', "Could not load current jobs.")
                return
            }

            let jobMap: { [key: string]: Job } = {}
            response.getJobsMap().forEach(function (value, key) {
                jobMap[key] = value
            })
            store.commit('updateJobData', jobMap)
            store.commit('updateTotalJobs')
        })
    }

    SubmitCreateForm(formulaData: CreateFormulaRequest.AsObject) {
        let self = this;
        let createFormulaRequest = new CreateFormulaRequest();
        createFormulaRequest.setName(formulaData.name);
        createFormulaRequest.setNumber(formulaData.number);
        createFormulaRequest.setNotes(formulaData.notes);
        createFormulaRequest.setJobsList(formulaData.jobsList);

        let basesList: Base[] = []
        formulaData.basesList.forEach(function (item, index) {
            let newBase = new Base();
            newBase.setType(item.type)
            newBase.setName(item.name)
            newBase.setAmount(item.amount)

            basesList.push(newBase)
        });

        createFormulaRequest.setBasesList(basesList);

        let colorantsList: Colorant[] = []
        formulaData.colorantsList.forEach(function (item, index) {
            let newColorant = new Colorant();
            newColorant.setType(item.type)
            newColorant.setName(item.name)
            newColorant.setAmount(item.amount)

            colorantsList.push(newColorant)
        })
        createFormulaRequest.setColorantsList(colorantsList);

        let metadata = { 'Authorization': 'Bearer ' + Vue.cookies.get('token') }
        this.client.createFormula(createFormulaRequest, metadata, function (err, response) {
            if (err) {
                console.log(err)
                // If formula already exists return helpful error
                if (err.code == 6) {
                    store.commit('displaySnackBar', "Could not create formula. Please make sure formula name is unique.")
                    return
                }
                store.commit('displaySnackBar', "Could not create formula.")
                return
            }
            store.commit("hideCreateFormulaModal")
            self.LoadFormulaData();
            self.LoadJobData();
            (app.$refs.createFormulaForm as HTMLFormElement).clearForm();
        })
    }
}
