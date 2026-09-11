import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgUpdateParams } from "./types/paktoken/tokenization/v1/tx";
import { MsgCreateProject } from "./types/paktoken/tokenization/v1/tx";
import { MsgUpdateProject } from "./types/paktoken/tokenization/v1/tx";
import { MsgDeleteProject } from "./types/paktoken/tokenization/v1/tx";
import { MsgBuyTokens } from "./types/paktoken/tokenization/v1/tx";
import { MsgDistributeRevenue } from "./types/paktoken/tokenization/v1/tx";
import { MsgCreateHolding } from "./types/paktoken/tokenization/v1/tx";
import { MsgUpdateHolding } from "./types/paktoken/tokenization/v1/tx";
import { MsgDeleteHolding } from "./types/paktoken/tokenization/v1/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/paktoken.tokenization.v1.MsgUpdateParams", MsgUpdateParams],
    ["/paktoken.tokenization.v1.MsgCreateProject", MsgCreateProject],
    ["/paktoken.tokenization.v1.MsgUpdateProject", MsgUpdateProject],
    ["/paktoken.tokenization.v1.MsgDeleteProject", MsgDeleteProject],
    ["/paktoken.tokenization.v1.MsgBuyTokens", MsgBuyTokens],
    ["/paktoken.tokenization.v1.MsgDistributeRevenue", MsgDistributeRevenue],
    ["/paktoken.tokenization.v1.MsgCreateHolding", MsgCreateHolding],
    ["/paktoken.tokenization.v1.MsgUpdateHolding", MsgUpdateHolding],
    ["/paktoken.tokenization.v1.MsgDeleteHolding", MsgDeleteHolding],
    
];

export { msgTypes }