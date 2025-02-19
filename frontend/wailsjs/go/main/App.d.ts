// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {main} from '../models';

export function AlertDialog(arg1:string,arg2:string):Promise<void>;

export function ConfirmDialog(arg1:string,arg2:string):Promise<string>;

export function GetSACNConfig():Promise<main.SACNConfig>;

export function LoadFile():Promise<string>;

export function SaveFile(arg1:string):Promise<void>;

export function SetCalibrationPoints(arg1:{[key: string]: main.CalibrationPoint}):Promise<void>;

export function SetFixtures(arg1:{[key: string]: main.Fixture}):Promise<void>;

export function SetMouseForAllFixtures(arg1:number,arg2:number):Promise<void>;

export function SetPanTiltForFixture(arg1:string,arg2:number,arg3:number):Promise<void>;

export function SetSACNConfig(arg1:main.SACNConfig):Promise<void>;

export function TypeExporter(arg1:main.CalibrationPoint,arg2:main.CalibratedCalibrationPoint,arg3:main.Fixture,arg4:main.SACNConfig,arg5:main.DMXData,arg6:main.Point):Promise<void>;
