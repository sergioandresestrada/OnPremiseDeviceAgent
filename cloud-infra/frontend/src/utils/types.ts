export declare module Jobs {

    export interface Message {
        "@hasStringResource": string;
        "#text": string;
    }

    export interface Status {
        Message: Message;
    }

    export interface Link {
        "@method": string;
        "@rel": string;
        "@uri": string;
    }

    export interface Links {
        Link: Link[];
    }

    export interface Job {
        Job_ID: string;
        Status: Status;
        Links: Links;
    }

    export interface Link2 {
        "@method": string;
        "@rel": string;
        "@uri": string;
    }

    export interface Links2 {
        Link: Link2[];
    }

    export interface Jobs {
        Version: string;
        Date: Date;
        Job: Job[];
        Links: Links2;
    }

    export interface RootObject {
        Jobs: Jobs;
    }

}

export declare module Identification {

    export interface InstalledRAM {
        Value: number;
        Units: string;
    }

    export interface InstalledHDD {
        Value: number;
        Units: string;
    }

    export interface DeviceType {
        Name: string;
        Version: string;
        Namespace: string;
    }

    export interface Fields {
        ModelName: string;
        ModelNumber: string;
        PartNumber: string;
        Manufacturer: string;
        SerialNumber: string;
        FriendlyName: string;
        DeviceRegion: string;
        FwReleaseName: string;
        FwReleaseDate: Date;
        InstalledRAM: InstalledRAM;
        InstalledHDD: InstalledHDD;
        DeviceType: DeviceType;
        FwInstallationDate: Date;
    }

    export interface P1 {
        X: number;
        Y: number;
        Z: number;
    }

    export interface P2 {
        X: number;
        Y: number;
        Z: number;
    }

    export interface MaxPlatform {
        P1: P1;
        P2: P2;
        Units: string;
    }

    export interface P12 {
        X: number;
        Y: number;
        Z: number;
    }

    export interface P22 {
        X: number;
        Y: number;
        Z: number;
    }

    export interface UsablePlatform {
        P1: P12;
        P2: P22;
        Units: string;
    }

    export interface PlatformAxes {
        XAxis: string;
        YAxis: string;
        ZAxis: string;
    }

    export interface BuildPlatform {
        MaxPlatform: MaxPlatform;
        UsablePlatform: UsablePlatform;
        PlatformAxes: PlatformAxes;
    }

    export interface BuildPlatforms {
        BuildPlatform: BuildPlatform;
    }

    export interface Namespaces {
        Namespace: any;
    }

    export interface SupportedExtension {
        Name: string;
        Version: string;
        Namespace: string;
        Namespaces: Namespaces;
    }

    export interface SupportedExtensions {
        SupportedExtension: SupportedExtension[];
    }

    export interface Namespaces2 {
        Namespace: string;
    }

    export interface RequiredExtension {
        Name: string;
        Version: string;
        Namespace: string;
        Namespaces: Namespaces2;
    }

    export interface RequiredExtensions {
        RequiredExtension: RequiredExtension;
    }

    export interface Namespaces3 {
        Namespace: string;
    }

    export interface ContentFormat {
        Name: string;
        Version: string;
        Namespace: string;
        SupportedExtensions: SupportedExtensions;
        RequiredExtensions: RequiredExtensions;
        Namespaces: Namespaces3;
    }

    export interface ContentFormats {
        ContentFormat: ContentFormat[];
    }

    export interface PrinterProperties {
        BuildPlatforms: BuildPlatforms;
        ContentFormats: ContentFormats;
        ColorSupported: string;
        JobAutoDelete: string;
    }

    export interface Name {
        "@hasStringResource": string;
        "#text": string;
    }

    export interface P13 {
        X: number;
        Y: number;
        Z: number;
    }

    export interface P23 {
        X: number;
        Y: number;
        Z: number;
    }

    export interface UsablePlatform2 {
        P1: P13;
        P2: P23;
        Units: string;
    }

    export interface PlatformAxes2 {
        XAxis: string;
        YAxis: string;
        ZAxis: string;
    }

    export interface BuildPlatform2 {
        UsablePlatform: UsablePlatform2;
        PlatformAxes: PlatformAxes2;
    }

    export interface Link {
        "@method": string;
        "@rel": string;
        "@uri": string;
    }

    export interface Links {
        Link: Link[];
    }

    export interface Material {
        ID: string;
        Name: Name;
        Default: boolean;
        BuildPlatform: BuildPlatform2;
        Links: Links;
    }

    export interface Materials {
        Material: Material[];
    }

    export interface Identification {
        Version: string;
        Date: Date;
        Fields: Fields;
        PrinterProperties: PrinterProperties;
        Materials: Materials;
    }

    export interface RootObject {
        Identification: Identification;
    }

}

