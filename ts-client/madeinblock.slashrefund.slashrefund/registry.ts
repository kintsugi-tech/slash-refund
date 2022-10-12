import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgWithdraw } from "./types/slashrefund/tx";
import { MsgDeposit } from "./types/slashrefund/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/madeinblock.slashrefund.slashrefund.MsgWithdraw", MsgWithdraw],
    ["/madeinblock.slashrefund.slashrefund.MsgDeposit", MsgDeposit],
    
];

export { msgTypes }