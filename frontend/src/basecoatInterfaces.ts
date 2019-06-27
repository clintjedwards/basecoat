import { Formula, Job } from "./basecoat_pb"

export interface LoginInfo {
    username: string;
    password: string;
}

export interface FormulaMap {
    [key: string]: Formula;
}

export interface JobMap {
    [key: string]: Job;
}

export interface colorantType {
    imageURL: string
    userMessage: string
}

export interface colorantTypeMap {
    [key: string]: colorantType;
}
