package main

type Provider func() ([]map[string]string, error)
