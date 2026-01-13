/* tslint:disable */
/* eslint-disable */
export function decrypt(data: EncryptedData, privkey: string, encoding_type: EncodingType): string;
export function sha256(data: string): string;
export function verify(data: string, sig: string, pubkey: string): boolean;
export function encrypt(data: string, pubkey: string, encoding_type: EncodingType): EncryptedData;
export function decrypt_aes(data: AesEncryptedData, key: string, encoding_type: EncodingType): string;
export function init_panic_hook(): void;
export function pedersen_reveal(commitment: Commitment, value: bigint, tag: string): boolean;
export function key_to_did(pubkey: string): Did;
export function ecdh(privkey: string, pubkey: string): string;
export function sign(data: string, privkey: string): string;
export function pedersen_commit(value: bigint, tag: string): Commitment;
export function encrypt_aes(data: string, key: string, encoding_type: EncodingType): AesEncryptedData;
export function create_keypair(): KeyPair;
export function did_to_key(did: Did): string;
export enum EncodingType {
  UTF8 = 0,
  HEX = 1,
  BASE64 = 2,
}
export class AesEncryptedData {
  free(): void;
  static from_object(object: any): AesEncryptedData;
  constructor(data: string, nonce: string);
  to_object(): any;
  readonly data: string;
  readonly nonce: string;
}
export class Commitment {
  free(): void;
  static from_object(object: any): Commitment;
  constructor(commitment: string, secret_blinding_factor: string);
  to_object(): any;
  commitment: string;
  secret_blinding_factor: string;
}
export class Did {
  free(): void;
  static from_object(object: any): Did;
  constructor(id: string);
  to_object(): any;
  id: string;
}
export class EncryptedData {
  free(): void;
  static from_object(object: any): EncryptedData;
  constructor(pubkey_from: string, pubkey_to: string, data: string, nonce: string);
  to_object(): any;
  readonly pubkey_from: string;
  readonly data: string;
  readonly nonce: string;
  readonly pubkey_to: string;
}
export class KeyPair {
  free(): void;
  static from_object(object: any): KeyPair;
  constructor(pubkey: string, privkey: string);
  to_object(): any;
  pubkey: string;
  privkey: string;
}
export class Tx {
  free(): void;
  static from_object(object: any): Tx;
  constructor(coin: string, from: string, to: string, amount: string, fee: string, data: string);
  to_object(): any;
  amount: string;
  to: string;
  fee: string;
  coin: string;
  data: string;
  from: string;
}
