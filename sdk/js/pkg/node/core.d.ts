/* tslint:disable */
/* eslint-disable */
export function create_keypair(): KeyPair;
export function key_to_did(pubkey: string): Did;
export function did_to_key(did: Did): string;
export function encrypt(data: string, pubkey: string): EncryptedData;
export function decrypt(data: EncryptedData, privkey: string): string;
export function sign(data: string, privkey: string): string;
export function verify(data: string, sig: string, pubkey: string): boolean;
export function sha256(data: string): string;
export function ecdh(privkey: string, pubkey: string): string;
export class Did {
  free(): void;
  constructor(id: string);
  to_object(): any;
  static from_object(object: any): Did;
  id: string;
}
export class EncryptedData {
  free(): void;
  constructor(pubkey_from: string, pubkey_to: string, data: string, nonce: string);
  to_object(): any;
  static from_object(object: any): EncryptedData;
  readonly pubkey_from: string;
  readonly pubkey_to: string;
  readonly data: string;
  readonly nonce: string;
}
export class KeyPair {
  free(): void;
  constructor(pubkey: string, privkey: string);
  to_object(): any;
  static from_object(object: any): KeyPair;
  pubkey: string;
  privkey: string;
}
export class Tx {
  free(): void;
  constructor(coin: string, from: string, to: string, amount: string, fee: string, data: string);
  to_object(): any;
  static from_object(object: any): Tx;
  coin: string;
  from: string;
  to: string;
  amount: string;
  fee: string;
  data: string;
}
