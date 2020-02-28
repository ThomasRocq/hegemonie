// Copyright (C) 2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
#ifndef HEGEMONIE_CMD_MAPGEN__BINDING_H
#define HEGEMONIE_CMD_MAPGEN__BINDING_H

#include <stdint.h>
#include <stdlib.h>
#include <assert.h>
#include <errno.h>

#include <igraph/igraph.h>

typedef struct vertex_s {
  int index;
  double x, y;
} vertex_t;

typedef struct vertex_array_s {
  uint32_t len;
  vertex_t *tab;
} vertex_array_t;

typedef struct edge_array_s {
  igraph_vector_t vector;
  uint32_t len;
} edge_array_t;


vertex_array_t *vertex_array_create(uint32_t nb) {
  vertex_array_t *va = malloc(sizeof(*va));
  if (!va)
    return NULL;

  va->tab = calloc(nb, sizeof(vertex_t));
  if (!va->tab) {
    free(va);
    return NULL;
  }

  va->len = nb;
  return va;
}

void vertex_array_destroy(vertex_array_t *va) {
  if (va) {
    if (va->tab)
      free(va->tab);
    free(va);
  }
}

void vertex_array_set(vertex_array_t *va, uint32_t index, vertex_t v) {
  assert(va != NULL);
  assert(index < va->len);
  va->tab[index] = v;
  va->tab[index].index = index;
}

double vertex_array_x(vertex_array_t *va, uint32_t index) {
  assert(va != NULL);
  assert(index < va->len);
  return va->tab[index].x;
}

double vertex_array_y(vertex_array_t *va, uint32_t index) {
  assert(va != NULL);
  assert(index < va->len);
  return va->tab[index].y;
}

edge_array_t *edge_array_create(uint32_t nb) {
  static const edge_array_t model = {IGRAPH_VECTOR_NULL, 0};
  edge_array_t *p = malloc(sizeof(*p));
  *p = model;
  igraph_vector_init(&p->vector, 0);
  igraph_vector_reserve(&p->vector, nb);
  return p;
}

void edge_array_destroy(edge_array_t *ea) {
  if (!ea)
    return;
  igraph_vector_destroy(&ea->vector);
  free(ea);
}

void edge_array_add(edge_array_t *ea, uint32_t from, uint32_t to) {
  assert(ea != NULL);
  igraph_vector_push_back(&ea->vector, from);
  igraph_vector_push_back(&ea->vector, to);
  ea->len++;
}

double clamp(double d, double lo, double hi) {
  if (d < lo) return lo;
  if (d > hi) return hi;
  return d;
}

int igraph_fdp(vertex_array_t *vertices, edge_array_t *edges,
    double x, double y, uint32_t niter) {
  igraph_t graph = {};
  igraph_matrix_t coord = {};
  igraph_vector_t minx = {}, maxx = {}, miny = {}, maxy = {};

  if (!vertices || vertices->len <= 0)
    return EAGAIN;
  if (!edges || edges->len <= 0)
    return EAGAIN;

  printf("V %u E %u\n", vertices->len, edges->len);
  igraph_create(&graph, &edges->vector, vertices->len, 1 /* directed */);
  igraph_matrix_init(&coord, vertices->len, 2);

  igraph_vector_init(&minx, 0);
  igraph_vector_init(&miny, 0);
  igraph_vector_init(&maxx, 0);
  igraph_vector_init(&maxy, 0);
  igraph_vector_reserve(&minx, vertices->len);
  igraph_vector_reserve(&miny, vertices->len);
  igraph_vector_reserve(&maxx, vertices->len);
  igraph_vector_reserve(&maxy, vertices->len);

  for (uint32_t i = 0; i < vertices->len; i++) {
    igraph_vector_push_back(&minx, 0.0);
    igraph_vector_push_back(&miny, 0.0);
    igraph_vector_push_back(&maxx, x);
    igraph_vector_push_back(&maxy, y);
    MATRIX(coord, i, 0) = clamp(vertices->tab[i].x, 0, x);
    MATRIX(coord, i, 1) = clamp(vertices->tab[i].y, 0, y);
  }
  printf("matrix rows %lu cols %lu\n", igraph_matrix_nrow(&coord), igraph_matrix_ncol(&coord));
  printf("graph vertices %d\n", igraph_vcount(&graph));

  int rc = igraph_layout_fruchterman_reingold(
      &graph, &coord,
      1 /*use_seed*/,
      niter /*niter*/,
      x /*start_temp*/,
      IGRAPH_LAYOUT_NOGRID /*grid*/,
      NULL /*weight*/,
      &minx /*minx*/,
      &maxx /*maxx*/,
      &miny /*miny*/,
      &maxy /*maxy*/);

  if (rc == 0) {
    assert(igraph_matrix_nrow(&coord) == vertices->len);
    for (int i = 0; i < igraph_matrix_nrow(&coord); i++) {
      vertices->tab[i].x = MATRIX(coord, i, 0);
      vertices->tab[i].y = MATRIX(coord, i, 1);
    }
  }

  igraph_vector_destroy(&minx);
  igraph_vector_destroy(&maxx);
  igraph_vector_destroy(&miny);
  igraph_vector_destroy(&maxy);
  igraph_matrix_destroy(&coord);
  igraph_destroy(&graph);
  return rc;
}

#endif  /* HEGEMONIE_CMD_MAPGEN__BINDING_H */
