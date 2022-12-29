import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgDeposit } from "./types/slashrefund/tx";
import { MsgClaim } from "./types/slashrefund/tx";
import { MsgWithdraw } from "./types/slashrefund/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/madeinblock.slashrefund.slashrefund.MsgDeposit", MsgDeposit],
    ["/madeinblock.slashrefund.slashrefund.MsgClaim", MsgClaim],
    ["/madeinblock.slashrefund.slashrefund.MsgWithdraw", MsgWithdraw],
    
];

export { msgTypes }