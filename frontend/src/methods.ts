import Vue from 'vue'
import VueCookies from 'vue-cookies'
import store from './store';
import * as bcInterface from './basecoatInterfaces'
import { app } from './index'

import {
    ListFormulasRequest,
    Formula, Job,
    ListJobsRequest,
    CreateAPITokenRequest,
    CreateFormulaRequest,
    Colorant,
    Base
} from "./basecoat_pb"
import { BasecoatClient } from './BasecoatServiceClientPb';

Vue.use(VueCookies)

// VerifyLogin verifies user login state
export function VerifyLogin(client: BasecoatClient) {
    if (!Vue.cookies.isKey('username') || !Vue.cookies.isKey('token')) {
        store.commit('updateLoginState', false)
        return
    }

    store.commit('updateUsername', Vue.cookies.get('username'))
    store.commit('updateLoginState', true)
    LoadFormulaData(client);
    LoadJobData(client);
}

// HandleLogin checks the backend for successful login and sets appropriate cookie
export function HandleLogin(client: BasecoatClient, loginInfo: bcInterface.LoginInfo) {
    store.commit('updateLoginIsLoading', true)

    let tokenRequest = new CreateAPITokenRequest();
    tokenRequest.setUser(loginInfo.username);
    tokenRequest.setPassword(loginInfo.password);
    tokenRequest.setDuration(10368000);
    client.createAPIToken(tokenRequest, {}, function (err, response) {
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
        LoadFormulaData(client)
        LoadJobData(client)
    })
}

// HandleLogout clears cookies and brings back up the login modal
export function HandleLogout(client: BasecoatClient) {
    Vue.cookies.remove('username')
    Vue.cookies.remove('token')
    VerifyLogin(client)
}

export function LoadFormulaData(client: BasecoatClient) {
    let listFormulasRequest = new ListFormulasRequest();
    let metadata = { 'Authorization': 'Bearer ' + Vue.cookies.get('token') }
    client.listFormulas(listFormulasRequest, metadata, function (err, response) {
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

export function LoadJobData(client: BasecoatClient) {
    let listJobsRequest = new ListJobsRequest();
    let metadata = { 'Authorization': 'Bearer ' + Vue.cookies.get('token') }
    client.listJobs(listJobsRequest, metadata, function (err, response) {
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

export function SubmitCreateForm(client: BasecoatClient, formulaData: CreateFormulaRequest.AsObject) {
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
    client.createFormula(createFormulaRequest, metadata, function (err, response) {
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
        LoadFormulaData(client);
        LoadJobData(client);
        (app.$refs.createFormulaForm as HTMLFormElement).clearForm();
    })
}


